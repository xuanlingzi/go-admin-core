package search

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// FromQueryTag tag标记
	FromQueryTag = "search"
	// Mysql 数据库标识
	Mysql = "mysql"
	// Postgres 数据库标识
	Postgres = "postgres"
)

// ResolveSearchQuery 解析
/**
 * 	exact / iexact 等于
 * 	contains / icontains 包含
 *	gt / gte 大于 / 大于等于
 *	lt / lte 小于 / 小于等于
 *	startswith / istartswith 以…起始
 *	endswith / iendswith 以…结束
 *	in
 *	isnull
 *  order 排序		e.g. order[key]=desc     order[key]=asc
 */
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	var tag string
	var ok bool
	var t *resolveSearchTag
	for i := 0; i < qType.NumField(); i++ {
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(FromQueryTag)
		if !ok {
			//递归调用
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), condition)
			continue
		}
		switch tag {
		case "-":
			continue
		}
		t = makeTag(tag)
		if qValue.Field(i).IsZero() {
			continue
		}
		//解析
		if t.Type != "order" && t.Custom != "" {
			condition.SetWhere(t.Custom, []interface{}{qValue.Field(i).Interface()})
		} else {
			column := fmt.Sprintf("`%s`.`%s`", t.Table, t.Column)
			if t.Func != "" {
				column = fmt.Sprintf("%s(`%s`.`%s`)", t.Func, t.Table, t.Column)
			}
			switch t.Type {
			case "left":
				//左关联
				join := condition.SetJoinOn(t.Type, fmt.Sprintf(
					"left join `%s` on `%s`.`%s` = `%s`.`%s`",
					t.Join,
					t.Join,
					t.On[0],
					t.Table,
					t.On[1],
				))
				ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
			case "exact", "iexact":
				condition.SetWhere(fmt.Sprintf("%s = ?", column), []interface{}{qValue.Field(i).Interface()})
			case "contains", "icontains":
				//fixme mysql不支持ilike
				if driver == Postgres && t.Type == "icontains" {
					condition.SetWhere(fmt.Sprintf("%s ilike ?", column), []interface{}{"%" + qValue.Field(i).String() + "%"})
				} else {
					condition.SetWhere(fmt.Sprintf("%s like ?", column), []interface{}{"%" + qValue.Field(i).String() + "%"})
				}
			case "gt":
				condition.SetWhere(fmt.Sprintf("%s > ?", column), []interface{}{qValue.Field(i).Interface()})
			case "gte":
				condition.SetWhere(fmt.Sprintf("%s >= ?", column), []interface{}{qValue.Field(i).Interface()})
			case "lt":
				condition.SetWhere(fmt.Sprintf("%s < ?", column), []interface{}{qValue.Field(i).Interface()})
			case "lte":
				condition.SetWhere(fmt.Sprintf("%s <= ?", column), []interface{}{qValue.Field(i).Interface()})
			case "startswith", "istartswith":
				if driver == Postgres && t.Type == "istartswith" {
					condition.SetWhere(fmt.Sprintf("%s ilike ?", column), []interface{}{qValue.Field(i).String() + "%"})
				} else {
					condition.SetWhere(fmt.Sprintf("%s like ?", column), []interface{}{qValue.Field(i).String() + "%"})
				}
			case "endswith", "iendswith":
				if driver == Postgres && t.Type == "iendswith" {
					condition.SetWhere(fmt.Sprintf("%s ilike ?", column), []interface{}{"%" + qValue.Field(i).String()})
				} else {
					condition.SetWhere(fmt.Sprintf("%s like ?", column), []interface{}{"%" + qValue.Field(i).String()})
				}
			case "in":
				condition.SetWhere(fmt.Sprintf("%s in (?)", column), []interface{}{qValue.Field(i).Interface()})
			case "isnull":
				if !(qValue.Field(i).IsZero() && qValue.Field(i).IsNil()) {
					condition.SetWhere(fmt.Sprintf("%s isnull", column), make([]interface{}, 0))
				}
			case "order":
				switch strings.ToLower(qValue.Field(i).String()) {
				case "desc", "asc":
					if t.Custom != "" {
						condition.SetOrder(fmt.Sprintf(t.Custom+" %s", qValue.Field(i).String()))
					} else {
						condition.SetOrder(fmt.Sprintf("%s %s", column, qValue.Field(i).String()))
					}
				}
			}
		}
	}
}
