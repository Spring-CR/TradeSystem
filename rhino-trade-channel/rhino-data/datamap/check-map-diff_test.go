package datamap_test

import (
	"fmt"
	"rhino-data/datamap"
	"testing"
)

func TestCheckDiff(t *testing.T) {
    tests := []struct {
        name           string
        mNew           map[string][]map[string]interface{}
        mOld           map[string][]map[string]interface{}
        wantAddKeys    []string
        wantChgKeys    []string
        wantDelKeys    []string
        wantErr        bool
    }{
        {
            name: "无变化",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "score": 95.5}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "score": 95.5}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "新增记录",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}},
                "key3": {{"name": "Charlie", "age": 35}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}},
            },
            wantAddKeys: []string{"key3"},
            wantChgKeys: []string{},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "删除记录",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{},
            wantDelKeys: []string{"key2"},
            wantErr:     false,
        },
        {
            name: "修改记录",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 26}}, // 年龄从25改为26
                "key2": {{"name": "Bob", "age": 30}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{"key1"},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "新增和删除同时存在",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key3": {{"name": "Charlie", "age": 35}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}},
            },
            wantAddKeys: []string{"key3"},
            wantChgKeys: []string{},
            wantDelKeys: []string{"key2"},
            wantErr:     false,
        },
        {
            name: "修改和删除同时存在",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 26}}, // 修改
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
                "key2": {{"name": "Bob", "age": 30}}, // 删除
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{"key1"},
            wantDelKeys: []string{"key2"},
            wantErr:     false,
        },
        {
            name: "内部map新增字段",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25, "city": "Beijing"}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{"key1"},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "内部map删除字段",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice"}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"name": "Alice", "age": 25}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{"key1"},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "多种数据类型",
            mNew: map[string][]map[string]interface{}{
                "key1": {{
                    "name":   "Alice",
                    "age":    int32(25),
                    "score":  float32(95.5),
                    "active": true,
                    "code":   uint8('A'),
                }},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{
                    "name":   "Alice",
                    "age":    int32(25),
                    "score":  float32(95.5),
                    "active": true,
                    "code":   uint8('A'),
                }},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "类型不匹配错误",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"age": 25}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"age": "25"}}, // 类型不匹配：int vs string
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{},
            wantDelKeys: []string{},
            wantErr:     true,
        },
        {
            name: "空map处理",
            mNew: map[string][]map[string]interface{}{},
            mOld: map[string][]map[string]interface{}{},
            wantAddKeys: []string{},
            wantChgKeys: []string{},
            wantDelKeys: []string{},
            wantErr:     false,
        },
        {
            name: "nil值处理",
            mNew: map[string][]map[string]interface{}{
                "key1": {{"value": nil}},
            },
            mOld: map[string][]map[string]interface{}{
                "key1": {{"value": "something"}},
            },
            wantAddKeys: []string{},
            wantChgKeys: []string{"key1"},
            wantDelKeys: []string{},
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            addKeys, chgKeys, delKeys, err := datamap.CheckDiff(tt.mNew, tt.mOld)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("datamap.CheckDiff() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                if !equalStringSlices(addKeys, tt.wantAddKeys) {
                    t.Errorf("datamap.CheckDiff() addKeys = %v, want %v", addKeys, tt.wantAddKeys)
                }
                if !equalStringSlices(chgKeys, tt.wantChgKeys) {
                    t.Errorf("datamap.CheckDiff() chgKeys = %v, want %v", chgKeys, tt.wantChgKeys)
                }
                if !equalStringSlices(delKeys, tt.wantDelKeys) {
                    t.Errorf("datamap.CheckDiff() delKeys = %v, want %v", delKeys, tt.wantDelKeys)
                }
            }
        })
    }
}

// 辅助函数：比较两个字符串切片是否相等（顺序无关）
func equalStringSlices(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    
    aMap := make(map[string]bool)
    for _, v := range a {
        aMap[v] = true
    }
    
    for _, v := range b {
        if !aMap[v] {
            return false
        }
    }
    
    return true
}

// 测试多个键的情况
func TestCheckDiffMultipleKeys(t *testing.T) {
    mNew := map[string][]map[string]interface{}{
        "add1": {{"name": "Add1"}},
        "add2": {{"name": "Add2"}},
        "chg1": {{"name": "Changed1"}},
        "chg2": {{"name": "Changed2"}},
        "same": {{"name": "Same"}},
    }
    
    mOld := map[string][]map[string]interface{}{
        "chg1": {{"name": "Change1"}}, // 会被修改
        "chg2": {{"name": "Change2"}}, // 会被修改
        "same": {{"name": "Same"}},    // 不变
        "del1": {{"name": "Delete1"}}, // 会被删除
        "del2": {{"name": "Delete2"}}, // 会被删除
    }
    
    addKeys, chgKeys, delKeys, err := datamap.CheckDiff(mNew, mOld)
    if err != nil {
        t.Errorf("datamap.CheckDiff() unexpected error: %v", err)
    }
    
    // 检查新增的键
    expectedAdd := []string{"add1", "add2"}
    if !equalStringSlices(addKeys, expectedAdd) {
        t.Errorf("Expected added keys %v, got %v", expectedAdd, addKeys)
    }
    
    // 检查修改的键
    expectedChg := []string{"chg1", "chg2"}
    if !equalStringSlices(chgKeys, expectedChg) {
        t.Errorf("Expected changed keys %v, got %v", expectedChg, chgKeys)
    }
    
    // 检查删除的键
    expectedDel := []string{"del1", "del2"}
    if !equalStringSlices(delKeys, expectedDel) {
        t.Errorf("Expected deleted keys %v, got %v", expectedDel, delKeys)
    }
}

// 性能测试
func BenchmarkCheckDiff(b *testing.B) {
    // 创建测试数据
    mNew := make(map[string][]map[string]interface{})
    mOld := make(map[string][]map[string]interface{})
    
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("key%d", i)
        mNew[key] = []map[string]interface{}{{
            "name":  fmt.Sprintf("Name%d", i),
            "age":   i,
            "score": float64(i) * 1.5,
        }}
        
        if i%2 == 0 { // 一半的数据在old中也有，但部分值不同
            mOld[key] = []map[string]interface{}{{
                "name":  fmt.Sprintf("Name%d", i),
                "age":   i + 1, // 年龄不同
                "score": float64(i) * 1.5,
            }}
        }
    }
    
    // 添加一些只在old中的数据
    for i := 1000; i < 1100; i++ {
        key := fmt.Sprintf("key%d", i)
        mOld[key] = []map[string]interface{}{{
            "name":  fmt.Sprintf("Name%d", i),
            "age":   i,
        }}
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        datamap.CheckDiff(mNew, mOld)
    }
}

// 测试特定类型的比较
func TestCheckDiffSpecificTypes(t *testing.T) {
    tests := []struct {
        name        string
        newVal      interface{}
        oldVal      interface{}
        shouldError bool
        shouldChange bool
    }{
        {"int32相同", int32(100), int32(100), false, false},
        {"int32不同", int32(100), int32(200), false, true},
        {"uint8相同", uint8(100), uint8(100), false, false},
        {"uint8不同", uint8(100), uint8(200), false, true},
        {"float32相同", float32(1.5), float32(1.5), false, false},
        {"float32不同", float32(1.5), float32(2.5), false, true},
        {"类型不匹配", int32(100), "100", true, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mNew := map[string][]map[string]interface{}{
                "key1": {{"value": tt.newVal}},
            }
            mOld := map[string][]map[string]interface{}{
                "key1": {{"value": tt.oldVal}},
            }
            
            _, _, _, err := datamap.CheckDiff(mNew, mOld)
            
            if tt.shouldError && err == nil {
                t.Errorf("Expected error but got none")
            } else if !tt.shouldError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            
            if !tt.shouldError && err == nil {
                _, chgKeys, _, _ := datamap.CheckDiff(mNew, mOld)
                hasChange := len(chgKeys) > 0
                if hasChange != tt.shouldChange {
                    t.Errorf("Expected change=%v, got change=%v", tt.shouldChange, hasChange)
                }
            }
        })
    }
}
