题目：给你一个链表数组，每个链表都已经按升序排列。请你将所有链表合并到一个升序链表中，返回合并后的链表

## 示例
```
输入：lists = [[1,4,5],[1,3,4],[2,6]]
输出：[1,1,2,3,4,4,5,6]
解释：链表数组如下：
[
  1->4->5,
  1->3->4,
  2->6
]
将它们合并到一个有序链表中得到。
1->1->2->3->4->4->5->6
```

## 方法一: 暴力求解
遍历所有链表，将所有节点的值放到一个数组中，然后将这个数组排序，然后遍历所有元素进行排序。最后用排序好的数组放进链表中

## 方法二: 逐一比较
比较K个节点(每个链表的首节点)，获得最小值的节点把它放在最终有序链表的后面，直到所有链表遍历完毕

## 方法三: 用队列方法优化方法2

## 方法四：逐一两两合并链表

## 方法五：分治
将k个链表两两配对合并，第一轮k个链表被合成两k/2个链表，接着继续合并，k/4，k/8个链表


## 代码
```go
func mergeTwoLists(list1, list2 *ListNode) *ListNode {
    dummy := &ListNode{} // 用哨兵节点简化代码逻辑
    cur := dummy         // cur 指向新链表的末尾
    for list1 != nil && list2 != nil {
        if list1.Val < list2.Val {
            cur.Next = list1 // 把 list1 加到新链表中
            list1 = list1.Next
        } else { // 注：相等的情况加哪个节点都是可以的
            cur.Next = list2 // 把 list2 加到新链表中
            list2 = list2.Next
        }
        cur = cur.Next
    }
    // 拼接剩余链表
    if list1 != nil {
        cur.Next = list1
    } else {
        cur.Next = list2
    }
    return dummy.Next
}

func mergeKLists(lists []*ListNode) *ListNode {
    m := len(lists)
    if m == 0 { // 注意输入的 lists 可能是空的
        return nil
    }
    if m == 1 { // 无需合并，直接返回
        return lists[0]
    }
    left := mergeKLists(lists[:m/2]) // 合并左半部分
    right := mergeKLists(lists[m/2:]) // 合并右半部分
    return mergeTwoLists(left, right) // 最后把左半和右半合并
}
```