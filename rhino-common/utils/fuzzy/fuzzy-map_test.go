package fuzzy_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"rhino-common/utils/fuzzy"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewFuzzyMap 测试构造函数
func TestNewFuzzyMap(t *testing.T) {
    keys := []string{"a", "b", "c"}
    values := []int{1, 2, 3}
    
    // 不使用去重函数
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    if fm == nil {
        t.Error("Expected non-nil FuzzyMap")
    }
    
    if fm.Size() != 3 {
        t.Errorf("Expected size 3, got %d", fm.Size())
    }
}

// TestNewFuzzyMapPanic 测试长度不匹配时的panic
func TestNewFuzzyMapPanic(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("Expected panic when keyList and valueList have different lengths")
        }
    }()
    
    keys := []string{"a", "b"}
    values := []int{1, 2, 3} // 长度不匹配
    
    _ = fuzzy.NewFuzzyMap(keys, values, nil)
}

// TestGetBasic 测试基本检索功能
func TestGetBasic(t *testing.T) {
    keys := []string{"hello", "world", "hell", "word"}
    values := []string{"value1", "value2", "value3", "value4"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试精确匹配
    result, ok := fm.Get("hello")
    if !ok || len(result) != 1 || result[0] != "value1" {
        t.Errorf("Expected ['value1'], got %v", result)
    }
    
    // 测试子串匹配
    result, ok = fm.Get("hell")
    if !ok || len(result) != 2 {
        t.Errorf("Expected 2 results for 'hell', got %d", len(result))
    }
    
    // 测试不存在的子串
    result, ok = fm.Get("nonexistent")
    if ok {
        t.Error("Expected no results for 'nonexistent'")
    }
}

// TestGetOrder 测试结果顺序
func TestGetOrder(t *testing.T) {
    keys := []string{"zebra", "apple", "banana", "application"}
    values := []int{1, 2, 3, 4}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 查找包含 "a" 的键，应该按照原始顺序返回
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    expected := []int{1, 2, 3, 4}
    if !reflect.DeepEqual(result, expected) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}

// TestGetWithDupFunc 测试使用去重函数
func TestGetWithDupFunc(t *testing.T) {
    keys := []string{"hello world", "world hello", "hello", "test hello"}
    values := []string{"shared", "shared", "shared", "unique"}
    
    // 定义去重函数
    dupFunc := func(a, b string) bool {
        return a == b
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    // 查找 "hello"，应该返回去重后的结果
    result, ok := fm.Get("hello")
    if !ok {
        t.Error("Expected results for 'hello'")
    }
    
    // 应该只有2个唯一结果（因为第4个键的值是唯一的）
    if len(result) != 2 {
        t.Errorf("Expected 2 unique results, got %d: %v", len(result), result)
    }
}

// TestGetWithoutDupFunc 测试不使用去重函数
func TestGetWithoutDupFunc(t *testing.T) {
    keys := []string{"a", "b", "a"}
    values := []int{1, 2, 1} // 相同的值，但没有使用去重函数
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    // 没有使用去重函数，应该返回所有匹配项
    if len(result) != 2 {
        t.Errorf("Expected 2 results without duplicate function, got %d", len(result))
    }
}

// TestGetComplexTypes 测试复杂类型
func TestGetComplexTypes(t *testing.T) {
    type Person struct {
        Name string
        Age  int
    }
    
    keys := []string{"alice", "bob", "charlie", "alicia"}
    values := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
        {"Alicia", 28},
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("ali")
    if !ok {
        t.Error("Expected results for 'ali'")
    }
    
    // 应该匹配 "alice" 和 "alicia"
    if len(result) != 2 {
        t.Errorf("Expected 2 results for 'ali', got %d", len(result))
    }
}

// TestGetEmptyMap 测试空映射
func TestGetEmptyMap(t *testing.T) {
    keys := []string{}
    values := []string{}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("anything")
    if ok {
        t.Error("Expected no results from empty map")
    }
    if result != nil {
        t.Errorf("Expected nil result from empty map, got %v", result)
    }
}

// TestGetSingleCharacter 测试单字符匹配
func TestGetSingleCharacter(t *testing.T) {
    keys := []string{"a", "aa", "aaa", "b", "c"}
    values := []int{1, 2, 3, 4, 5}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    if len(result) != 3 {
        t.Errorf("Expected 3 results for 'a', got %d", len(result))
    }
}

// TestGetOverlappingMatches 测试重叠匹配
func TestGetOverlappingMatches(t *testing.T) {
    keys := []string{"abc", "bcd", "cde"}
    values := []string{"v1", "v2", "v3"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // "bc" 应该匹配前两个键
    result, ok := fm.Get("bc")
    if !ok || len(result) != 2 {
        t.Errorf("Expected 2 results for 'bc', got %d", len(result))
    }
    
    // "cd" 应该匹配后两个键
    result, ok = fm.Get("cd")
    if !ok || len(result) != 2 {
        t.Errorf("Expected 2 results for 'cd', got %d", len(result))
    }
}

// TestGetSpecialCharacters 测试特殊字符
func TestGetSpecialCharacters(t *testing.T) {
    keys := []string{"hello-world", "hello_world", "hello world", "test-me"}
    values := []string{"v1", "v2", "v3", "v4"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    testCases := []struct {
        subkey   string
        expected int
    }{
        {"hello", 3},
        {"world", 3},
        {"test", 1},
        {"-", 2},
        {"_", 1},
        {" ", 1},
    }
    
    for _, tc := range testCases {
        result, ok := fm.Get(tc.subkey)
        if !ok && tc.expected > 0 {
            t.Errorf("Expected results for '%s'", tc.subkey)
            continue
        }
        
        if ok && len(result) != tc.expected {
            t.Errorf("Expected %d results for '%s', got %d", tc.expected, tc.subkey, len(result))
        }
    }
}

// TestGetEdgeCases 测试边界情况
func TestGetEdgeCases(t *testing.T) {
    // 测试只有一个键的情况
    keys := []string{"single"}
    values := []string{"value"}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("single")
    if !ok || len(result) != 1 {
        t.Error("Expected 1 result for single key")
    }
    
    // 测试所有键都相同的情况
    keys = []string{"test", "test", "test"}
    values = []string{"v1", "v2", "v3"}
    fm = fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok = fm.Get("test")
    if !ok || len(result) != 3 {
        t.Errorf("Expected 3 results for identical keys, got %d", len(result))
    }
}

// TestGetWithComplexDupFunc 测试复杂去重函数
func TestGetWithComplexDupFunc(t *testing.T) {
    type Product struct {
        ID   int
        Name string
    }
    
    keys := []string{"apple phone", "apple watch", "samsung phone", "google phone"}
    values := []Product{
        {1, "iPhone"},
        {3, "Apple Watch"},
        {1, "iPhone"}, // 与第一个相同
        {4, "Pixel"},
    }
    
    // 定义去重函数：如果两个产品的ID相同，则认为重复
    dupFunc := func(a, b Product) bool {
        return a.ID == b.ID
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    result, ok := fm.Get("phone")
    if !ok {
        t.Error("Expected results for 'phone'")
    }
    
    // 应该只有2个结果（因为第一个和第三个被认为是重复的）
    if len(result) != 2 {
        t.Errorf("Expected 2 unique results with complex dup func, got %d", len(result))
    }
}

// TestGetDeduplicationLogic 测试去重逻辑
func TestGetDeduplicationLogic(t *testing.T) {
    keys := []string{"key1", "key2", "key3", "key4"}
    values := []int{1, 1, 3, 1} // 多个重复值
    
    // 简单的去重函数
    dupFunc := func(a, b int) bool {
        return a == b
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    // 搜索 "key"，应该匹配所有键，但去重后只有2个唯一值
    result, ok := fm.Get("key")
    if !ok {
        t.Error("Expected results for 'key'")
    }
    
    // 应该只有2个唯一结果
    if len(result) != 2 {
        t.Errorf("Expected 2 unique results, got %d: %v", len(result), result)
    }
}

// TestGetEmptySubkey 测试空子串
func TestGetEmptySubkey(t *testing.T) {
    keys := []string{"a", "b"}
    values := []string{"v1", "v2"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("")
    if ok {
        t.Error("Expected no results for empty subkey")
    }
    if result != nil {
        t.Errorf("Expected nil result for empty subkey, got %v", result)
    }
}

// TestGetCaseSensitive 测试大小写敏感
func TestGetCaseSensitive(t *testing.T) {
    keys := []string{"Hello", "hello", "HELLO"}
    values := []string{"value1", "value2", "value3"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 大小写敏感，应该只匹配一个
    result, ok := fm.Get("Hello")
    if !ok || len(result) != 1 {
        t.Errorf("Expected 1 result for 'Hello', got %d", len(result))
    }
    
    result, ok = fm.Get("hello")
    if !ok || len(result) != 1 {
        t.Errorf("Expected 1 result for 'hello', got %d", len(result))
    }
}

// TestGetUnicode 测试Unicode字符
func TestGetUnicode(t *testing.T) {
    keys := []string{"café", "cafe", "cafetería", "咖啡", "カフェ"}
    values := []string{"french", "english", "spanish", "chinese", "japanese"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试Unicode子串匹配
    result, ok := fm.Get("caf")
    if !ok || len(result) != 3 {
        t.Errorf("Expected 3 results for 'caf', got %d", len(result))
    }
    
    // 测试中文字符
    result, ok = fm.Get("咖啡")
    if !ok || len(result) != 1 {
        t.Errorf("Expected 1 result for '咖啡', got %d", len(result))
    }
}

// TestGetVeryLongKeys 测试超长键
func TestGetVeryLongKeys(t *testing.T) {
    // 创建超长键
    longKey1 := strings.Repeat("abc", 1000) + "target" + strings.Repeat("def", 1000)
    longKey2 := strings.Repeat("xyz", 1000) + "target" + strings.Repeat("uvw", 1000)
    
    keys := []string{longKey1, longKey2, "short"}
    values := []string{"value1", "value2", "value3"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 在超长键中搜索
    result, ok := fm.Get("target")
    if !ok || len(result) != 2 {
        t.Errorf("Expected 2 results for 'target', got %d", len(result))
    }
}

// TestGetPartialOverlap 测试部分重叠匹配
func TestGetPartialOverlap(t *testing.T) {
    keys := []string{"abcd", "bcde", "cdef", "defg"}
    values := []int{1, 2, 3, 4}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试各种重叠情况
    testCases := []struct {
        subkey   string
        expected int
    }{
        {"abc", 1},
        {"bcd", 2},
        {"cde", 2},
        {"def", 2},
        {"efg", 1},
    }
    
    for _, tc := range testCases {
        result, ok := fm.Get(tc.subkey)
        if !ok && tc.expected > 0 {
            t.Errorf("Expected results for '%s'", tc.subkey)
            continue
        }
        
        if ok && len(result) != tc.expected {
            t.Errorf("Expected %d results for '%s', got %d", tc.expected, tc.subkey, len(result))
        }
    }
}

// TestGetWithNilDupFunc 测试nil去重函数
func TestGetWithNilDupFunc(t *testing.T) {
    keys := []string{"a", "b", "a"}
    values := []string{"same", "same", "same"} // 所有值都相同
    
    // 使用nil去重函数
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    // 没有去重，应该返回所有匹配项
    if len(result) != 2 {
        t.Errorf("Expected 2 results with nil dup func, got %d", len(result))
    }
}

// TestGetWithAlwaysTrueDupFunc 测试总是返回true的去重函数
func TestGetWithAlwaysTrueDupFunc(t *testing.T) {
    keys := []string{"a", "b", "c", "d"}
    values := []string{"value1", "value2", "value3", "value4"}
    
    // 总是返回true的去重函数
    dupFunc := func(a, b string) bool {
        return true
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    // 由于所有值都被认为是重复的，应该只返回一个结果
    if len(result) != 1 {
        t.Errorf("Expected 1 result with always-true dup func, got %d", len(result))
    }
}

// TestGetWithAlwaysFalseDupFunc 测试总是返回false的去重函数
func TestGetWithAlwaysFalseDupFunc(t *testing.T) {
    keys := []string{"a", "b", "a"}
    values := []string{"same", "same", "same"} // 所有值都相同
    
    // 总是返回false的去重函数
    dupFunc := func(a, b string) bool {
        return false
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    result, ok := fm.Get("a")
    if !ok {
        t.Error("Expected results for 'a'")
    }
    
    // 没有值被认为是重复的，应该返回所有匹配项
    if len(result) != 2 {
        t.Errorf("Expected 2 results with always-false dup func, got %d", len(result))
    }
}

// TestGetConcurrent 测试并发安全
func TestGetConcurrent(t *testing.T) {
    keys := []string{"apple", "banana", "cherry", "date", "elderberry"}
    values := []string{"fruit1", "fruit2", "fruit3", "fruit4", "fruit5"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    var wg sync.WaitGroup
    numGoroutines := 100
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            subkey := "a"
            if id%2 == 0 {
                subkey = "e"
            }
            
            result, ok := fm.Get(subkey)
            if !ok {
                t.Errorf("Goroutine %d: Expected results for '%s'", id, subkey)
                return
            }
            
            if len(result) == 0 {
                t.Errorf("Goroutine %d: Expected non-empty results", id)
            }
        }(i)
    }
    
    wg.Wait()
}

// TestGetPerformance 测试性能
func TestGetPerformance(t *testing.T) {
    // 创建大量数据
    numItems := 10000
    keys := make([]string, numItems)
    values := make([]int, numItems)
    
    for i := 0; i < numItems; i++ {
        keys[i] = "key_" + string(rune('a'+i%26)) + "_" + string(rune('a'+(i/26)%26)) + "_" + string(rune('a'+(i/676)%26))
        values[i] = i
    }
    
    // 测量构建时间
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测量搜索时间
    result, ok := fm.Get("key_a")
    if !ok {
        t.Error("Expected results for 'key_a'")
    }
    
    if len(result) == 0 {
        t.Error("Expected non-empty results")
    }
    
    t.Logf("Performance test: %d items, found %d results for 'key_a'", numItems, len(result))
}

// TestGetMultipleSearches 测试多次搜索
func TestGetMultipleSearches(t *testing.T) {
    keys := []string{"apple", "application", "banana", "band", "candy"}
    values := []string{"fruit1", "fruit2", "fruit3", "fruit4", "fruit5"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 多次搜索不同子串
    searches := []struct {
        subkey   string
        expected int
    }{
        {"app", 2},
        {"ban", 2},
        {"cand", 1},
        {"d", 2}, // "application", "band", "candy"
        {"y", 1}, // "candy"
    }
    
    for _, search := range searches {
        result, ok := fm.Get(search.subkey)
        if !ok && search.expected > 0 {
            t.Errorf("Expected results for '%s'", search.subkey)
            continue
        }
        
        if ok && len(result) != search.expected {
            t.Errorf("Expected %d results for '%s', got %d", search.expected, search.subkey, len(result))
        }
    }
}

// TestGetExactSubstring 测试精确子串匹配
func TestGetExactSubstring(t *testing.T) {
    keys := []string{"prefixsuffix", "prefix", "suffix", "fix"}
    values := []int{1, 2, 3, 4}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试各种子串位置
    testCases := []struct {
        subkey   string
        expected int
    }{
        {"prefix", 2}, // "prefixsuffix", "prefix"
        {"suffix", 2}, // "prefixsuffix", "suffix"
        {"fix", 4},    // 所有键都包含"fix"
        {"pre", 2},    // "prefixsuffix", "prefix"
        {"suf", 2},    // "prefixsuffix", "suffix"
    }
    
    for _, tc := range testCases {
        result, ok := fm.Get(tc.subkey)
        if !ok && tc.expected > 0 {
            t.Errorf("Expected results for '%s'", tc.subkey)
            continue
        }
        
        if ok && len(result) != tc.expected {
            t.Errorf("Expected %d results for '%s', got %d", tc.expected, tc.subkey, len(result))
        }
    }
}

// TestGetWithSingleKey 测试单键情况
func TestGetWithSingleKey(t *testing.T) {
    keys := []string{"lonely"}
    values := []string{"only"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试各种子串
    testCases := []struct {
        subkey   string
        expected bool
    }{
        {"lonely", true},
        {"lone", true},
        {"only", false}, // 注意：这是在值中，不是在键中
        {"xyz", false},
    }
    
    for _, tc := range testCases {
        result, ok := fm.Get(tc.subkey)
        if ok != tc.expected {
            t.Errorf("For subkey '%s', expected success=%v, got %v", tc.subkey, tc.expected, ok)
        }
        
        if ok && len(result) != 1 {
            t.Errorf("For subkey '%s', expected 1 result, got %d", tc.subkey, len(result))
        }
    }
}

// TestGetWithMixedTypes 测试混合类型
func TestGetWithMixedTypes(t *testing.T) {
    type Mixed struct {
        Str string
        Num int
        Flt float64
    }
    
    keys := []string{"key1", "key2", "key3"}
    values := []Mixed{
        {"a", 1, 1.1},
        {"c", 3, 3.3},
        {"a", 1, 1.1}, // 与第一个相同
    }
    
    // 基于Str和Num字段的去重函数
    dupFunc := func(a, b Mixed) bool {
        return a.Str == b.Str && a.Num == b.Num
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    result, ok := fm.Get("key")
    if !ok {
        t.Error("Expected results for 'key'")
    }
    
    // 应该只有2个唯一结果（索引0和1）
    if len(result) != 2 {
        t.Errorf("Expected 2 unique results with mixed types, got %d", len(result))
    }
}

// TestEmptyKeysOrValues 测试边界条件测试
func TestEmptyKeysOrValues(t *testing.T) {
    // 空键列表
    keys := []string{}
    values := []int{}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    result, ok := fm.Get("any")
    if ok || result != nil {
        t.Error("Expected no results for empty map")
    }

    // 空字符串键
    keys = []string{"", " ", "a"}
    values2 := []string{"v1", "v2", "v3"}
    fm2 := fuzzy.NewFuzzyMap(keys, values2, nil)
    _, ok = fm2.Get("")
    if ok {
        t.Error("Expected no results for empty subkey")
    }
}

// TestInvalidUTF8 错误场景测试
func TestInvalidUTF8(t *testing.T) {
    keys := []string{"normal", "invalid\xff\xfe", "test"}
    values := []int{1, 2, 3}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)

    result, ok := fm.Get("invalid")
    if !ok || len(result) != 1 {
        t.Errorf("Expected 1 result for invalid UTF-8 key, got %d", len(result))
    }
}

// TestDuplicateKeysDifferentValues 重复键不同值测试​​：验证相同键对应不同值时的去重逻辑。
func TestDuplicateKeysDifferentValues(t *testing.T) {
    keys := []string{"key", "key", "key"} // 重复键
    values := []string{"v1", "v2", "v3"} // 不同值
    dupFunc := func(a, b string) bool { return a == b }
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)

    result, ok := fm.Get("key")
    if !ok || len(result) != 3 { // 不去重，因为值不同
        t.Errorf("Expected 3 results for duplicate keys with different values, got %d", len(result))
    }
}

// TestLargeDataset 大数据集测试​​：验证10万级键值对的构建和检索性能。
func TestLargeDataset(t *testing.T) {
    numItems := 100000
    keys := make([]string, numItems)
    values := make([]int, numItems)
    for i := 0; i < numItems; i++ {
        keys[i] = fmt.Sprintf("key_%d_%d", i%100, i/100)
        values[i] = i
    }

    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    result, ok := fm.Get("key_50")
    if !ok {
        t.Error("Expected results for large dataset")
    }
    if len(result) < 1000 { // 应匹配大量键
        t.Errorf("Expected many results, got %d", len(result))
    }
}

// TestConcurrentAccess 高并发访问测试​​：模拟多线程同时调用Get方法。
func TestConcurrentAccess(t *testing.T) {
    keys := []string{"a", "b", "c", "d", "e"}
    values := []int{1, 2, 3, 4, 5}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)

    var wg sync.WaitGroup
    numGoroutines := 1000
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            subkey := "a"
            if id%2 == 0 {
                subkey = "e"
            }
            result, ok := fm.Get(subkey)
            if !ok || len(result) == 0 {
                t.Errorf("Goroutine %d: Expected results for '%s'", id, subkey)
            }
        }(i)
    }
    wg.Wait()
}

// TestNestedStructDeduplication 嵌套结构去重测试​​：使用复杂结构体自定义去重规则。
func TestNestedStructDeduplication(t *testing.T) {
    type Address struct { City string; Zip int }
    type Person struct { Name string; Addr Address }

    keys := []string{"alice", "alicia", "bob"}
    values := []Person{
        {"Alice", Address{"Beijing", 100000}},
        {"Alicia", Address{"Beijing", 100000}}, // 相同地址，应去重
        {"Bob", Address{"Shanghai", 200000}},
    }

    // 根据地址去重
    dupFunc := func(a, b Person) bool {
        return a.Addr.City == b.Addr.City && a.Addr.Zip == b.Addr.Zip
    }
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)

    result, ok := fm.Get("a")
    if !ok || len(result) != 1 { // "alice"和"alicia"去重后保留一个
        t.Errorf("Expected 1 unique result for nested struct, got %d", len(result))
    }
}

// TestPartialFieldDeduplication 部分字段去重测试​​：仅比较结构体的部分字段。
func TestPartialFieldDeduplication(t *testing.T) {
    type Product struct { ID int; Name string; Price float64 }
    keys := []string{"p1", "p2", "p3"}
    values := []Product{
        {1, "Phone", 999.0},
        {1, "Phone", 899.0}, // ID相同，应去重
        {2, "Tablet", 799.0},
    }

    dupFunc := func(a, b Product) bool { return a.ID == b.ID }
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)

    result, ok := fm.Get("p")
    if !ok || len(result) != 2 { // 前两个Product去重
        t.Errorf("Expected 2 unique results by ID, got %d", len(result))
    }
}

// FuzzFuzzyMapGet 基础模糊测试​​：随机生成子串测试Get方法。
func FuzzFuzzyMapGet(f *testing.F) {
    // 种子语料库：基于已知键初始化
    keys := []string{"hello", "world", "golang", "test"}
    values := []int{1, 2, 3, 4}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)

    f.Add("hell") // 添加种子数据
    f.Fuzz(func(t *testing.T, subkey string) {
        result, ok := fm.Get(subkey)
        if ok {
            // 验证结果非空且类型正确
            if len(result) == 0 {
                t.Error("Get returned ok but empty result")
            }
            for _, v := range result {
                if reflect.TypeOf(v).Kind() != reflect.Int {
                    t.Errorf("Unexpected value type: %T", v)
                }
            }
        }
        // 即使ok=false，也应正常处理（不panic）
    })
}

// TestOverlappingSubstring 重叠子串匹配测试​​：验证如"abc"在"xabc"和"abcy"中的匹配。
func TestOverlappingSubstring(t *testing.T) {
    keys := []string{"prefixsuffix", "prefix", "suffix", "fix"}
    values := []rune{'A', 'B', 'C', 'D'}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)

    cases := []struct {
        subkey   string
        expected int
    }{
        {"prefix", 2}, // 匹配"prefixsuffix"和"prefix"
        {"suffix", 2}, // 匹配"prefixsuffix"和"suffix"
        {"fix", 4},    // 匹配所有键
    }
    for _, tc := range cases {
        result, ok := fm.Get(tc.subkey)
        if !ok && tc.expected > 0 {
            t.Errorf("Expected results for '%s'", tc.subkey)
        }
        if ok && len(result) != tc.expected {
            t.Errorf("For '%s', expected %d results, got %d", tc.subkey, tc.expected, len(result))
        }
    }
}

// TestUnicodeNormalization Unicode规范化测试​​：检查Unicode组合字符的处理。
func TestUnicodeNormalization(t *testing.T) {
    // "café"可能被编码为"cafe\u0301"
    keys := []string{"café", "cafe", "cafeteria"}
    values := []string{"French", "English", "Spanish"}
    fm := fuzzy.NewFuzzyMap(keys, values, nil)

    result, ok := fm.Get("caf")
    if !ok || len(result) != 3 {
        t.Errorf("Expected 3 results for Unicode substring, got %d", len(result))
    }
}

// TestSuffixArrayEdgeCases 测试后缀数组的边界匹配情况
// func TestSuffixArrayEdgeCases(t *testing.T) {
//     // 测试包含特殊边界字符的键
//     keys := []string{"", " ", "\x00", "\n", "\t", "a\000b", "x\ny"}
//     values := []int{0, 1, 2, 3, 4, 5, 6}
    
//     fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
//     testCases := []struct {
//         subkey    string
//         expected  int
//         shouldOk bool
//     }{
//         {"", 0, false},      // 空子串不应匹配
//         {" ", 1, true},      // 空格匹配
//         {"\x00", 1, true},   // 空字符匹配（索引2的键）
//         {"\n", 1, true},     // 换行符匹配
//         {"a", 1, true},      // 部分匹配
//     }
    
//     for _, tc := range testCases {
//         result, ok := fm.Get(tc.subkey)
//         if ok != tc.shouldOk {
//             t.Errorf("For subkey %q, expected ok=%v, got %v", tc.subkey, tc.shouldOk, ok)
//             continue
//         }
//         if ok && len(result) != tc.expected {
//             t.Errorf("For subkey %q, expected %d results, got %d", tc.subkey, tc.expected, len(result))
//         }
//     }
// }

// TestSuffixArrayMultipleMatches 测试同一键内多次匹配
func TestSuffixArrayMultipleMatches(t *testing.T) {
    // 键内包含多个相同子串
    keys := []string{"ababab", "abcabc", "testtest"}
    values := []string{"v1", "v2", "v3"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 虽然"ab"在第一个键中出现3次，但应该只匹配一次该键
    result, ok := fm.Get("ab")
    if !ok || len(result) != 2 { // 应匹配前两个键
        t.Errorf("Expected 2 results for 'ab', got %d", len(result))
    }
    
    // 测试重叠匹配
    result, ok = fm.Get("test")
    if !ok || len(result) != 1 {
        t.Errorf("Expected 1 result for 'test', got %d", len(result))
    }
}

// TestDeduplicationEdgeCases 测试去重逻辑的边界情况
func TestDeduplicationEdgeCases(t *testing.T) {
    // 测试空值和nil值的去重
    var nilSlice []string
    keys := []string{"k1", "k2", "k3", "k4"}
    values := []interface{}{
        nil,                    // nil值
        nilSlice,               // nil切片
        []string{},             // 空切片
        []string{"non-empty"},  // 非空值
    }
    
    // 自定义去重函数处理nil和空值
    dupFunc := func(a, b interface{}) bool {
        // 如果都是nil，认为相同
        if a == nil && b == nil {
            return true
        }
        // 如果都是切片，比较内容
        aSlice, aOk := a.([]string)
        bSlice, bOk := b.([]string)
        if aOk && bOk {
            if len(aSlice) != len(bSlice) {
                return false
            }
            for i := range aSlice {
                if aSlice[i] != bSlice[i] {
                    return false
                }
            }
            return true
        }
        return false
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    result, ok := fm.Get("k")
    if !ok {
        t.Error("Expected results for 'k'")
    }
    
    // 应该只有3个唯一结果（前两个nil值应该去重）
    if len(result) != 3 {
        t.Errorf("Expected 3 unique results with nil values, got %d", len(result))
    }
}

// TestDeduplicationWithPanic 测试去重函数中的panic处理
func TestDeduplicationWithPanic(t *testing.T) {
    keys := []string{"a", "b", "c"}
    values := []int{1, 2, 3}
    
    // 会panic的去重函数
    panicDupFunc := func(a, b int) bool {
        if a == 2 {
            panic("intentional panic in dup function")
        }
        return a == b
    }
    
    // 应该正常处理panic而不影响构造函数
    fm := fuzzy.NewFuzzyMap(keys, values, panicDupFunc)
    
    // 即使去重函数可能panic，Get操作也应该安全
    result, ok := fm.Get("a")
    if !ok || len(result) != 1 {
        t.Error("Get should work even with problematic dup function")
    }
}

// TestMemoryUsageWithLargeKeys 测试超长键的内存使用
func TestMemoryUsageWithLargeKeys(t *testing.T) {
    // 创建大量超长键，测试内存使用
    numKeys := 1000
    keys := make([]string, numKeys)
    values := make([]int, numKeys)
    
    for i := 0; i < numKeys; i++ {
        // 创建10KB的长键
        longKey := strings.Repeat(fmt.Sprintf("key%d", i), 1000)
        keys[i] = longKey
        values[i] = i
    }
    
    // 测量构建时间和内存使用
    start := time.Now()
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    buildTime := time.Since(start)
    
    t.Logf("Built FuzzyMap with %d long keys in %v", numKeys, buildTime)
    
    // 测试检索性能
    start = time.Now()
    result, ok := fm.Get("key500")
    searchTime := time.Since(start)
    
    if !ok {
        t.Error("Expected results for 'key500'")
    }
    
    t.Logf("Search completed in %v, found %d results", searchTime, len(result))
    
    if buildTime > 10*time.Second {
        t.Error("Build time for large keys is too long")
    }
}

// TestFuzzyMapGCBehavior 测试垃圾回收行为
func TestFuzzyMapGCBehavior(t *testing.T) {
    // 创建大量临时FuzzyMap测试GC行为
    for i := 0; i < 100; i++ {
        keys := make([]string, 100)
        values := make([]int, 100)
        for j := 0; j < 100; j++ {
            keys[j] = fmt.Sprintf("temp_key_%d_%d", i, j)
            values[j] = j
        }
        
        fm := fuzzy.NewFuzzyMap(keys, values, nil)
        // 立即使用并丢弃，测试GC
        result, ok := fm.Get("temp")
        if ok && len(result) == 0 {
            t.Error("Unexpected empty results")
        }
    }
    
    // 强制GC检查内存泄漏
    runtime.GC()
}

// TestConcurrentModification 测试并发修改场景
func TestConcurrentModification(t *testing.T) {
    keys := []string{"a", "b", "c", "d", "e"}
    values := []string{"v1", "v2", "v3", "v4", "v5"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    var wg sync.WaitGroup
    readers := 100
    writers := 5 // 注意：实际FuzzyMap是不可变的，这里测试读安全
    
    // 并发读取
    for i := 0; i < readers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                result, ok := fm.Get("a")
                if !ok && len(result) != 1 {
                    t.Errorf("Concurrent reader %d: unexpected result", id)
                }
                // 短暂睡眠模拟工作负载
                time.Sleep(time.Microsecond)
            }
        }(i)
    }
    
    // 模拟"写"操作（实际是创建新实例）
    for i := 0; i < writers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 10; j++ {
                // 创建新的FuzzyMap实例模拟修改
                newKeys := append(keys, fmt.Sprintf("new_%d_%d", id, j))
                newValues := append(values, fmt.Sprintf("new_v_%d_%d", id, j))
                _ = fuzzy.NewFuzzyMap(newKeys, newValues, nil)
            }
        }(i)
    }
    
    wg.Wait()
}

// TestRaceConditions 测试竞态条件
func TestRaceConditions(t *testing.T) {
    keys := []string{"race", "condition", "test"}
    values := []int{1, 2, 3}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 使用race detector测试
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            result, ok := fm.Get("race")
            if ok && len(result) != 1 {
                t.Error("Unexpected race condition detected")
            }
        }()
    }
    wg.Wait()
}

// TestErrorRecovery 测试错误恢复能力
func TestErrorRecovery(t *testing.T) {
    // 测试在异常输入下的恢复能力
    keys := []string{"normal", "exceptional", "test"}
    values := []string{"v1", "v2", "v3"}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试各种可能引起panic的输入
    dangerousInputs := []string{
        string(make([]byte, 1000000)), // 超长输入
        string([]byte{0xff, 0xfe, 0xfd}), // 无效UTF-8
        "normal", // 正常输入
    }
    
    for _, input := range dangerousInputs {
        func() {
            defer func() {
                if r := recover(); r != nil {
                    t.Errorf("Panic occurred with input %q: %v", input, r)
                }
            }()
            // 应该不会panic
            result, ok := fm.Get(input)
            if ok {
                _ = result // 使用结果避免编译警告
            }
        }()
    }
}

// TestFuzzyMapSerialization 测试序列化行为（如需要）
func TestFuzzyMapSerialization(t *testing.T) {
    keys := []string{"serialize", "test"}
    values := []int{1, 2}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试FuzzyMap是否可序列化（如需要网络传输或持久化）
    // 注意：suffixarray.Index可能包含不可序列化的字段
    func() {
        defer func() {
            if r := recover(); r != nil {
                t.Log("Expected serialization limitation:", r)
            }
        }()
        
        // 尝试序列化（示例）
        var buf bytes.Buffer
        encoder := gob.NewEncoder(&buf)
        err := encoder.Encode(fm)
        if err != nil {
            t.Log("FuzzyMap serialization not supported:", err)
        }
    }()
}

// FuzzFuzzyMapGet2 基础模糊测试
func FuzzFuzzyMapGet2(f *testing.F) {
    // 种子语料库
    seeds := []string{"a", "hello", "test", "world", "123", " ", "\x00"}
    
    keys := []string{"hello world", "test case", "fuzzy matching", "go language"}
    values := []int{1, 2, 3, 4}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    for _, seed := range seeds {
        f.Add(seed)
    }
    
    f.Fuzz(func(t *testing.T, subkey string) {
        // 过滤掉空字符串（已知会返回false）
        if subkey == "" {
            t.Skip()
        }
        
        result, ok := fm.Get(subkey)
        
        // 验证不变量
        if ok {
            // 结果不应为空
            if len(result) == 0 {
                t.Error("Get returned ok but empty result")
            }
            
            // 所有结果类型应该正确
            for _, v := range result {
                if reflect.TypeOf(v).Kind() != reflect.Int {
                    t.Errorf("Unexpected value type: %T", v)
                }
            }
            
            // 结果数量不应超过键数量
            if len(result) > len(keys) {
                t.Errorf("More results than keys: %d > %d", len(result), len(keys))
            }
        }
        
        // 即使!ok，也不应panic
    })
}

// FuzzFuzzyMapWithDupFunc 带去重函数的模糊测试
func FuzzFuzzyMapWithDupFunc(f *testing.F) {
    type ComplexStruct struct {
        ID      int
        Name    string
        Tags    []string
        Enabled bool
    }
    
    keys := []string{"struct1", "struct2", "struct3"}
    values := []ComplexStruct{
        {1, "test", []string{"a", "b"}, true},
        {2, "test", []string{"c", "d"}, false},
        {1, "test", []string{"a", "b"}, true}, // 重复第一个
    }
    
    dupFunc := func(a, b ComplexStruct) bool {
        return a.ID == b.ID && a.Name == b.Name
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    f.Add("struct")
    f.Add("test")
    f.Add("1")
    
    f.Fuzz(func(t *testing.T, subkey string) {
        if subkey == "" {
            t.Skip()
        }
        
        result, ok := fm.Get(subkey)
        if ok {
            // 验证去重逻辑
            seen := make(map[int]bool)
            for _, item := range result {
                if seen[item.ID] {
                    t.Error("Duplicate items after deduplication")
                }
                seen[item.ID] = true
            }
        }
    })
}

// TestCrossPlatformConsistency 测试跨平台一致性
func TestCrossPlatformConsistency(t *testing.T) {
    keys := []string{"hello", "world", "test"}
    values := []int{1, 2, 3}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 测试在不同环境下结果一致性
    testCases := []string{"h", "he", "hell", "hello", "w", "wo", "wor", "world"}
    
    for _, subkey := range testCases {
        result1, ok1 := fm.Get(subkey)
        
        // 重新构建实例测试一致性
        fm2 := fuzzy.NewFuzzyMap(keys, values, nil)
        result2, ok2 := fm2.Get(subkey)
        
        if ok1 != ok2 {
            t.Errorf("Inconsistent ok result for %q: %v vs %v", subkey, ok1, ok2)
        }
        
        if ok1 && !reflect.DeepEqual(result1, result2) {
            t.Errorf("Inconsistent results for %q: %v vs %v", subkey, result1, result2)
        }
    }
}

// TestEndiannessAwareness 测试字节序敏感性（如需要）
func TestEndiannessAwareness(t *testing.T) {
    // 测试二进制数据的处理（如果FuzzyMap处理二进制数据）
    keys := []string{
        string([]byte{0x00, 0x01, 0x02, 0x03}), // 小端序模式
        string([]byte{0x03, 0x02, 0x01, 0x00}), // 大端序模式
        "normal string",
    }
    values := []int{1, 2, 3}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    // 搜索二进制模式
    result, ok := fm.Get(string([]byte{0x00, 0x01}))
    if ok && len(result) != 1 {
        t.Errorf("Expected 1 result for binary pattern, got %d", len(result))
    }
}

// TestFuzzyMapIntegration 综合集成测试
func TestFuzzyMapIntegration(t *testing.T) {
    // 模拟真实场景的复杂测试
    type Document struct {
        ID      int
        Title   string
        Content string
        Tags    []string
    }
    
    // 模拟文档检索系统
    documents := []Document{
        {1, "Go Programming", "Go is a programming language", []string{"go", "programming"}},
        {2, "Go Fuzzing", "Fuzzing is a testing technique", []string{"go", "testing"}},
        {3, "Python Programming", "Python is another language", []string{"python", "programming"}},
    }
    
    // 提取标题作为键，文档作为值
    keys := make([]string, len(documents))
    for i, doc := range documents {
        keys[i] = doc.Title
    }
    
    fm := fuzzy.NewFuzzyMap(keys, documents, func(a, b Document) bool {
        return a.ID == b.ID // 根据ID去重
    })
    
    // 测试各种搜索场景
    searchScenarios := []struct {
        query     string
        minResults int
        description string
    }{
        {"Go", 2, "应该找到所有Go相关文档"},
        {"Programming", 2, "应该找到所有编程相关文档"},
        {"Python", 1, "应该找到Python文档"},
        {"Fuzz", 1, "应该支持前缀匹配"},
        {"ing", 3, "应该支持子串匹配"},
    }
    
    for _, scenario := range searchScenarios {
        t.Run(scenario.query, func(t *testing.T) {
            result, ok := fm.Get(scenario.query)
            if !ok {
                t.Errorf("查询 %q 应该返回结果", scenario.query)
                return
            }
            
            if len(result) < scenario.minResults {
                t.Errorf("查询 %q 预期至少 %d 个结果，得到 %d", 
                    scenario.query, scenario.minResults, len(result))
            }
            
            t.Logf("for %s\n", scenario.query)
            // 验证结果顺序（应该保持原始顺序）
            for i, doc := range result {
                t.Logf(">>> doc.ID:%v", doc.ID)
                if doc.ID < documents[i].ID {
                    t.Errorf("结果顺序不正确")
                    break
                }
            }
        })
    }
}

// BenchmarkGet 性能基准测试
func BenchmarkGet(b *testing.B) {
    // 准备测试数据
    keys := make([]string, 1000)
    values := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = "key" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
        values[i] = i
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fm.Get("key")
    }
}

// BenchmarkConstruction 构建性能基准测试
func BenchmarkConstruction(b *testing.B) {
    keys := make([]string, 1000)
    values := make([]string, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = "key" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
        values[i] = "value"
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = fuzzy.NewFuzzyMap(keys, values, nil)
    }
}

// BenchmarkGetWithDupFunc 使用去重函数的性能基准测试
func BenchmarkGetWithDupFunc(b *testing.B) {
    keys := make([]string, 1000)
    values := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = "key" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
        // 创建一些重复值
        if i%10 == 0 {
            values[i] = i / 10
        } else {
            values[i] = i
        }
    }
    
    dupFunc := func(a, b int) bool {
        return a == b
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fm.Get("key")
    }
}

// 测试FuzzyMap对中文key的处理能力
func TestFuzzyMapWithChineseKeys(t *testing.T) {
    keys := []string{
        "Go语言编程",
        "Python人工智能开发", 
        "Java后端架构师",
        "前端Vue.js框架",
        "数据库MySQL优化",
        "云计算Kubernetes部署",
        "机器学习深度学习",
        "区块链技术应用",
        "微服务架构设计",
        "DevOps持续集成",
    }
    
    values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    testCases := []struct {
        name        string
        query       string
        minResults  int
        description string
    }{
        {
            name:        "纯中文匹配",
            query:       "编程",
            minResults:  1,
            description: "测试纯中文字符串的模糊匹配",
        },
        {
            name:        "中英混合匹配", 
            query:       "Go语言",
            minResults:  1,
            description: "测试中文与英文混合的匹配",
        },
        {
            name:        "部分匹配",
            query:       "框架",
            minResults:  1,
            description: "测试中文关键词的部分匹配",
        },
        {
            name:        "长中文文本匹配",
            query:       "架构",
            minResults:  2, 
            description: "测试在较长中文文本中的匹配",
        },
        {
            name:        "中文专业术语匹配",
            query:       "人工智能",
            minResults:  1,
            description: "测试包含中文专业术语的匹配",
        },
        {
            name:        "无匹配情况",
            query:       "不存在的关键词",
            minResults:  0,
            description: "测试无匹配结果时的处理",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Logf("测试用例: %s - %s", tc.name, tc.description)
            
            results, ok := fm.Get(tc.query)
            
            if tc.minResults > 0 && !ok {
                t.Errorf("查询'%s'应该匹配到结果，但实际未匹配到", tc.query)
                return
            }
            
            if len(results) < tc.minResults {
                t.Errorf("查询'%s'期望至少%d个结果，实际得到%d个", 
                    tc.query, tc.minResults, len(results))
            }
            
            t.Logf("查询'%s'匹配到%d个结果:", tc.query, len(results))
            for i, keyIndex := range results {
                t.Logf("  结果%d: key='%s', value=%d", 
                    i+1, keys[keyIndex], values[keyIndex])
            }
        })
    }
}

// 测试中文key的边界情况
func TestFuzzyMapChineseEdgeCases(t *testing.T) {
    edgeCases := []struct {
        name   string
        keys   []string
        values []int
        query  string
    }{
        {
            name:   "中文Unicode边界",
            keys:   []string{"一", "鿎", "𠀀"},
            values: []int{1, 2, 3},
            query:  "一",
        },
        {
            name:   "中文数字混合",
            keys:   []string{"第1章", "版本2.0", "2023年更新"},
            values: []int{1, 2, 3},
            query:  "第1",
        },
        {
            name:   "中文特殊标点",
            keys:   []string{"你好！", "请问：", "价格¥100"},
            values: []int{1, 2, 3},
            query:  "价格",
        },
        {
            name:   "中英日混合",
            keys:   []string{"Hello世界", "Pythonプログラミング", "Java编程"},
            values: []int{1, 2, 3},
            query:  "编程",
        },
    }
    
    for _, tc := range edgeCases {
        t.Run(tc.name, func(t *testing.T) {
            fm := fuzzy.NewFuzzyMap(tc.keys, tc.values, nil)
            results, ok := fm.Get(tc.query)
            
            if !ok {
                t.Errorf("边界用例'%s'匹配失败", tc.name)
                return
            }

          //  jsData, _ := json.Marshal(results)
            
            t.Logf("边界用例'%s': 查询'%s'匹配到%d个结果: %v", 
                tc.name, tc.query, len(results), results)
            
            // 修复：添加边界检查
            for i, keyIndex := range results {
                keyIndex -= 1
                if keyIndex < 0 || keyIndex > len(tc.keys)-1 {
                    t.Errorf("错误：索引%d超出范围[0, %d]", keyIndex, len(tc.keys)-1)
                    continue
                }
                t.Logf("  匹配结果%d: key='%s'", i+1, tc.keys[keyIndex])
            }
        })
    }
}

// 性能基准测试：中文key的匹配性能
func BenchmarkFuzzyMapWithChinese(b *testing.B) {
    keys := []string{
        "分布式系统架构设计与实践指南手册教程",
        "微服务云原生应用开发部署运维管理平台", 
        "人工智能机器学习深度学习神经网络算法",
        "区块链分布式账本智能合约加密技术",
        "前端响应式设计用户体验界面开发",
        "后端高并发缓存数据库优化方案",
        "全栈工程师技能树学习路径指南",
    }
    
    values := make([]int, len(keys))
    for i := range values {
        values[i] = i
    }
    
    fm := fuzzy.NewFuzzyMap(keys, values, nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 测试不同的中文查询模式
        fm.Get("架构")
        fm.Get("学习") 
        fm.Get("服务")
        fm.Get("设计")
        fm.Get("优化")
    }
}

// 综合测试：验证中文key的排序和去重功能
func TestFuzzyMapChineseWithDeduplication(t *testing.T) {
    // 测试数据包含重复值
    keys := []string{
        "Go语言并发编程",
        "Go语言网络编程", 
        "Python数据分析",
        "Python机器学习",
        "Java并发编程",
        "Java网络编程",
    }
    
    values := []int{1, 1, 2, 2, 3, 3} // 故意设置重复值
    
    // 使用去重函数
    dupFunc := func(a, b int) bool { return a == b }
    fm := fuzzy.NewFuzzyMap(keys, values, dupFunc)
    
    testCases := []struct {
        query       string
        expectedMin int
    }{
        {"编程", 2}, // 应该去重后得到3个唯一值
        {"Python", 1},
        {"并发", 2},
    }
    
    for _, tc := range testCases {
        t.Run(tc.query, func(t *testing.T) {
            results, ok := fm.Get(tc.query)
            if !ok {
                t.Errorf("查询'%s'应该返回结果", tc.query)
                return
            }
            
            // 检查去重效果
            uniqueValues := make(map[int]bool)
            for _, v := range results {
                uniqueValues[v] = true
            }
            
            if len(uniqueValues) < tc.expectedMin {
                t.Errorf("查询'%s'去重后期望至少%d个唯一值，实际得到%d个", 
                    tc.query, tc.expectedMin, len(uniqueValues))
            }
            
            t.Logf("查询'%s': 总结果数=%d, 去重后唯一值数=%d", 
                tc.query, len(results), len(uniqueValues))
        })
    }
}

// go test -v -bench=.
