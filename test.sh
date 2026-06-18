#!/bin/bash
# 图像 MCP 服务器集成测试脚本

BASE_URL="http://localhost:8080"
TEST_IMG="sample_waveform.png"
PASS=0
FAIL=0
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

pass() { echo -e "${GREEN}✅ PASS${NC}: $1"; PASS=$((PASS+1)); }
fail() { echo -e "${RED}❌ FAIL${NC}: $1"; FAIL=$((FAIL+1)); }

# 读取配置中的 Token（如果有）
if [ -f config.json ]; then
    AUTH_TOKEN=$(python3 -c "import json; print(json.load(open('config.json')).get('auth_token',''))" 2>/dev/null)
fi
if [ -n "$AUTH_TOKEN" ]; then
    AUTH_HEADER=(-H "Authorization: Bearer $AUTH_TOKEN")
else
    AUTH_HEADER=()
fi

echo "========================================="
echo "  图像 MCP 服务器 集成测试"
echo "========================================="

# 检查服务器是否在运行
echo ""
echo "--- 检查服务器状态 ---"
if curl -s -o /dev/null -w "%{http_code}" "${AUTH_HEADER[@]}" "$BASE_URL/upload" | grep -q "405\|400\|200"; then
    pass "服务器在运行"
else
    fail "服务器未启动，请先运行 ./img-mcp-server"; exit 1
fi
[ -n "$AUTH_TOKEN" ] && echo "  (已启用 Token 认证)"

# ======== analyze_image 测试 ========

# 测试1: 上传图片
echo ""
echo "--- 测试1: 上传图片 ---"
CACHE_ID1=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/upload" \
    -F "file=@$TEST_IMG" | python3 -c "import sys,json; print(json.load(sys.stdin)['cache_id'])" 2>/dev/null)
[ -n "$CACHE_ID1" ] && pass "上传成功" || fail "上传失败"

# 测试2: 同一张图再上传，cache_id 应不变
echo ""
echo "--- 测试2: 哈希去重 ---"
CACHE_ID2=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/upload" \
    -F "file=@$TEST_IMG" | python3 -c "import sys,json; print(json.load(sys.stdin)['cache_id'])" 2>/dev/null)
[ "$CACHE_ID1" = "$CACHE_ID2" ] && pass "同一图片 cache_id 一致" \
    || fail "同一图片 cache_id 不同"

# 测试3: 用 cache_id 分析
echo ""
echo "--- 测试3: analyze_image(cache_id) ---"
RESULT=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"analyze_image\",\"arguments\":{\"cache_id\":\"$CACHE_ID1\",\"prompt\":\"用一句话描述\"}}}" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['result']['content'][0]['text'])" 2>/dev/null)
[ -n "$RESULT" ] && pass "分析成功: ${RESULT:0:60}" || fail "分析返回空"

# 测试4: 无效 cache_id
echo ""
echo "--- 测试4: 无效 cache_id ---"
ERR=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"analyze_image","arguments":{"cache_id":"not-exist","prompt":"test"}}}' \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['result']['content'][0]['text'])" 2>/dev/null)
echo "$ERR" | grep -q "缓存" && pass "正确返回错误" || fail "错误信息不对: $ERR"

# 测试5: 缺参数
echo ""
echo "--- 测试5: 缺图片参数 ---"
ERR2=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"analyze_image","arguments":{"prompt":"test"}}}' \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['result']['content'][0]['text'])" 2>/dev/null)
echo "$ERR2" | grep -q "提供" && pass "正确返回错误" || fail "错误信息不对: $ERR2"

# ======== compare_images 测试 ========

# 准备第二张测试图
echo ""
echo "--- 准备对比测试图 ---"
python3 -c "from PIL import Image; Image.new('RGB', (100, 100), color=(255,0,0)).save('/tmp/test_red.png')" 2>/dev/null \
    && pass "对比测试图已生成" || fail "生成对比测试图失败"

# 测试6: 工具列表包含 compare_images
echo ""
echo "--- 测试6: tools/list 包含 compare_images ---"
TOOLS=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' \
    | python3 -c "import sys,json; print([t['name'] for t in json.load(sys.stdin)['result']['tools']])" 2>/dev/null)
echo "$TOOLS" | grep -q "compare_images" && pass "compare_images 工具已注册" \
    || fail "compare_images 工具未找到: $TOOLS"

# 测试7: 上传第二张图并对比
echo ""
echo "--- 测试7: compare_images ---"
RED_ID=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/upload" \
    -F "file=@/tmp/test_red.png" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['cache_id'])" 2>/dev/null)
CMP=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"compare_images\",\"arguments\":{\"cache_id_1\":\"$CACHE_ID1\",\"cache_id_2\":\"$RED_ID\",\"prompt\":\"对比两张图\"}}}" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['result']['content'][0]['text'])" 2>/dev/null)
[ -n "$CMP" ] && pass "对比成功: ${CMP:0:60}" || fail "对比返回空"

# 测试8: 缺第二张图参数
echo ""
echo "--- 测试8: 缺图2参数 ---"
ERR3=$(curl -s -X POST "${AUTH_HEADER[@]}" "$BASE_URL/mcp" -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"compare_images\",\"arguments\":{\"cache_id_1\":\"$CACHE_ID1\",\"prompt\":\"test\"}}}" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['result']['content'][0]['text'])" 2>/dev/null)
echo "$ERR3" | grep -q "提供" && pass "缺图2正确返回错误" || fail "错误信息不对: $ERR3"

echo ""
echo "========================================="
echo -e "  结果: ${GREEN}$PASS 通过${NC}, ${RED}$FAIL 失败${NC}"
echo "========================================="
[ "$FAIL" -eq 0 ] || exit 1
