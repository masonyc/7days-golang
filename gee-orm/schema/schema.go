package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

//Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string //eg. primary key, Not Null etc..
}

//Schema represents a table of database
type Schema struct {
	Model      interface{} //Entity Type
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field), //init map
	}

	for i := 0; i < modelType.NumField(); i++ {
		// NumField() 获取实例的字段的个数，然后通过下标获取到特定字段 p := modelType.Field(i)
		p := modelType.Field(i)
		//IsExported reports whether name starts with an upper-case letter.
		// .Anonymous check struct has a name
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
