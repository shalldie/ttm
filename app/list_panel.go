package app

import (
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type IListPanel interface {
	setFocus()
}

type ListPanel[T any] struct {
	*tview.Flex
	list             *tview.List       // 列表组件
	newItem          *tview.InputField // 新加项组件
	tipBar           *tview.TextView
	items            []*T              // 列表
	parent           IListPanel        // 上一个panel
	child            IListPanel        // 上一个panel
	activeItem       *T                // 新加项
	loadFromDB       func()            // 从db中获取数据
	addNewItem       func(text string) // 添加新项
	onSelectedChange func(item *T)     // 选择项改变
}

func newListPanel[T any](title string, newItemText string) *ListPanel[T] {
	// instance
	l := &ListPanel[T]{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
		list: tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true).
			SetSelectedStyle(
				tcell.Style{}.Background(tcell.ColorBlue),
			),
		newItem: makeLightTextInput(" + [" + newItemText + "] "),
		tipBar:  tview.NewTextView().SetText(" 新建：N ; 删除：D").SetTextColor(tcell.ColorYellow),
	}

	// 组件
	l.SetBorder(true).SetTitle(" " + title + " ")
	l.AddItem(l.list, 0, 1, true).AddItem(l.newItem, 1, 0, false).AddItem(l.tipBar, 1, 0, false)
	// 兼容 powerlevel10k
	l.list.SetBorderPadding(0, 0, 1, 1)
	l.newItem.SetBorderPadding(0, 0, 1, 1)

	// 事件 - list
	l.SetFocusFunc(func() {
		l.setFocus()
	})
	// SetSelectedFunc
	l.list.SetChangedFunc(func(i int, s1, s2 string, r rune) {
		l.activeItem = l.items[i]
		if l.onSelectedChange != nil {
			l.onSelectedChange(l.activeItem)
		}
	})

	// 事件 - newproject
	l.newItem.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			l.addNewItem(strings.TrimSpace(l.newItem.GetText()))
		case tcell.KeyEsc:
			l.setFocus()
		}
	})

	return l

}

// 重置数据、状态
func (l *ListPanel[T]) reset() {
	l.list.Clear()
	l.items = make([]*T, 0)
	l.activeItem = nil
	if l.loadFromDB != nil {
		l.loadFromDB()
	}
	l.newItem.SetText("")
}

// 设置焦点
func (l *ListPanel[T]) setFocus() {
	app.SetFocus(l)
}

// 处理快捷键
func (l *ListPanel[T]) handleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch unicode.ToLower(event.Rune()) {
	// 新建
	case 'n':
		app.SetFocus(l.newItem)
		return nil
	// 删除
	case 'd':
		app.SetFocus(l.newItem)
		return nil
	}

	// 向左
	if event.Key() == tcell.KeyLeft && l.parent != nil {
		l.parent.setFocus()
		return nil
	}
	// 向右
	if event.Key() == tcell.KeyRight && l.child != nil {
		l.child.setFocus()
		return nil
	}

	return event
}
