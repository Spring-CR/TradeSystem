package fuzzy

import (
	"fmt"
	"index/suffixarray"
	"log"
	"rhino-common/domain_error"
	"sort"
	"strings"

	"github.com/manucorporat/try"
)

// FuzzyMap 模糊映射类，支持通过子串检索值
type FuzzyMap[T any] struct {
	keys      []string
	values    []T
	sa        *suffixarray.Index
	allData   []byte
	keyRanges []keyRange
	valueMap  map[int]int // 用于去重的映射
}

type keyRange struct {
	start, end int
	index      int // 在原始 keys 和 values 中的索引
}

// NewFuzzyMap 创建新的 FuzzyMap 实例
// dupFunc 用于判断两个值是否重复，如果为 nil 则不去重
func NewFuzzyMap[T any](keyList []string, valueList []T, dupFunc func(a, b T) bool) *FuzzyMap[T] {
	if len(keyList) != len(valueList) {
		panic("keyList and valueList must have the same length")
	}

	// 构建所有键连接起来的大字符串
	var builder strings.Builder
	keyRanges := make([]keyRange, len(keyList))

	for i, key := range keyList {
		start := builder.Len()
		builder.WriteString(key)
		end := builder.Len()
		keyRanges[i] = keyRange{
			start: start,
			end:   end,
			index: i,
		}
	}

	allData := []byte(builder.String())
	sa := suffixarray.New(allData)

	// 构建 valueMap 用于去重
	valueMap := make(map[int]int)
	if dupFunc != nil {
		// 使用提供的去重函数
		for i := range valueList {
			if _, exists := valueMap[i]; !exists {
				valueMap[i] = i
				// 查找所有重复的值
				for j := i + 1; j < len(valueList); j++ {
					var dup bool
					try.This(func() {
						dup = dupFunc(valueList[i], valueList[j])
					}).Catch(func(err try.E) {
						log.Printf("error occur while run dupFunc! error:%v\n", err)
						de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while run dupFunc! error:%v", err))
						domain_error.ProcessSevereError(false, 0, de, nil, fmt.Sprintf("error occur while run dupFunc! error:%v\n", err))
					})
					if _, exists := valueMap[j]; !exists && dup {
						valueMap[j] = i
					}
				}
			}
		}
	} else {
		// 不去重，每个索引指向自己
		for i := range valueList {
			valueMap[i] = i
		}
	}

	return &FuzzyMap[T]{
		keys:      keyList,
		values:    valueList,
		sa:        sa,
		allData:   allData,
		keyRanges: keyRanges,
		valueMap:  valueMap,
	}
}

// Get 根据子串检索值，返回结果的顺序与输入的 valueList 保持一致
func (fm *FuzzyMap[T]) Get(subkey string) ([]T, bool) {
	if subkey == "" {
		return nil, false
	}

	// 使用后缀数组查找所有匹配位置
	indices := fm.sa.Lookup([]byte(subkey), -1)
	if len(indices) == 0 {
		return nil, false
	}

	// 找出所有匹配的键索引
	matchedIndices := make(map[int]bool)
	for _, pos := range indices {
		// 查找包含该位置的键
		for _, kr := range fm.keyRanges {
			if pos >= kr.start && pos+len(subkey) <= kr.end {
				matchedIndices[kr.index] = true
				break
			}
		}
	}

	if len(matchedIndices) == 0 {
		return nil, false
	}

	// 转换为原始值索引并去重
	uniqueValueIndices := make(map[int]bool)
	for idx := range matchedIndices {
		// 使用 valueMap 来获取唯一的代表索引
		representativeIdx := fm.valueMap[idx]
		uniqueValueIndices[representativeIdx] = true
	}

	// 将唯一的索引转换为切片并按原始顺序排序
	sortedIndices := make([]int, 0, len(uniqueValueIndices))
	for idx := range uniqueValueIndices {
		sortedIndices = append(sortedIndices, idx)
	}

	// 按原始输入顺序排序
	sort.Ints(sortedIndices)

	// 收集结果
	result := make([]T, 0, len(sortedIndices))
	for _, idx := range sortedIndices {
		result = append(result, fm.values[idx])
	}

	return result, true
}

// GetKeys 辅助方法：获取所有键
func (fm *FuzzyMap[T]) GetKeys() []string {
	return fm.keys
}

// Size 辅助方法：获取映射大小
func (fm *FuzzyMap[T]) Size() int {
	return len(fm.keys)
}

// Dispose 清理数据并释放内存，调用后实例将不可用
func (fm *FuzzyMap[T]) Dispose() {

	// 清空所有切片，帮助GC回收内存
	fm.keys = nil
	fm.values = nil
	fm.keyRanges = nil
	fm.allData = nil
	
	// 清空映射
	fm.valueMap = nil
	
	// 注意：suffixarray.Index 没有显式的关闭方法
	// 我们将其设为 nil 以便 GC 可以回收
	fm.sa = nil
}