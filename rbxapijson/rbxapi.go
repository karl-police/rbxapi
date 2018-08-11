// The rbxapijson package implements a codec for the Roblox API dump in JSON
// format.
package rbxapijson

import (
	"github.com/robloxapi/rbxapi"
)

type Root struct {
	Classes []*Class
	Enums   []*Enum
}

// GetClasses implements the rbxapi.Root interface.
func (root *Root) GetClasses() []rbxapi.Class {
	list := make([]rbxapi.Class, len(root.Classes))
	for i, class := range root.Classes {
		list[i] = class
	}
	return list
}

// GetClass implements the rbxapi.Root interface.
func (root *Root) GetClass(name string) rbxapi.Class {
	for _, class := range root.Classes {
		if class.Name == name {
			return class
		}
	}
	return nil
}

// GetEnums implements the rbxapi.Root interface.
func (root *Root) GetEnums() []rbxapi.Enum {
	list := make([]rbxapi.Enum, len(root.Enums))
	for i, enum := range root.Enums {
		list[i] = enum
	}
	return list
}

// GetEnum implements the rbxapi.Root interface.
func (root *Root) GetEnum(name string) rbxapi.Enum {
	for _, enum := range root.Enums {
		if enum.Name == name {
			return enum
		}
	}
	return nil
}

type Class struct {
	Name           string
	Superclass     string
	MemoryCategory string
	Members        []rbxapi.Member
	Tags           `json:",omitempty"`
}

// GetName implements the rbxapi.Class interface.
func (class *Class) GetName() string {
	return class.Name
}

// GetSuperclass implements the rbxapi.Class interface.
func (class *Class) GetSuperclass() string {
	return class.Superclass
}

// GetMembers implements the rbxapi.Class interface.
func (class *Class) GetMembers() []rbxapi.Member {
	list := make([]rbxapi.Member, len(class.Members))
	copy(list, class.Members)
	return list
}

// GetMember implements the rbxapi.Class interface.
func (class *Class) GetMember(name string) rbxapi.Member {
	for _, member := range class.Members {
		if member.GetName() == name {
			return member
		}
	}
	return nil
}

type Property struct {
	Name          string
	ValueType     Type
	Category      string
	ReadSecurity  string
	WriteSecurity string
	CanLoad       bool
	CanSave       bool
	Tags          `json:",omitempty"`
}

func (member *Property) GetMemberType() string     { return "Property" }
func (member *Property) GetName() string           { return member.Name }
func (member *Property) GetValueType() rbxapi.Type { return member.ValueType }
func (member *Property) GetSecurity() (read, write string) {
	return member.ReadSecurity, member.WriteSecurity
}

type Function struct {
	Name       string
	Parameters []Parameter
	ReturnType Type
	Security   string
	Tags       `json:",omitempty"`
}

func (member *Function) GetMemberType() string      { return "Function" }
func (member *Function) GetName() string            { return member.Name }
func (member *Function) GetReturnType() rbxapi.Type { return member.ReturnType }
func (member *Function) GetSecurity() string        { return member.Security }
func (member *Function) GetParameters() []rbxapi.Parameter {
	list := make([]rbxapi.Parameter, len(member.Parameters))
	for i, param := range member.Parameters {
		list[i] = param
	}
	return list
}

type Event struct {
	Name       string
	Parameters []Parameter
	Security   string
	Tags       `json:",omitempty"`
}

func (member *Event) GetMemberType() string { return "Event" }
func (member *Event) GetName() string       { return member.Name }
func (member *Event) GetSecurity() string   { return member.Security }
func (member *Event) GetParameters() []rbxapi.Parameter {
	list := make([]rbxapi.Parameter, len(member.Parameters))
	for i, param := range member.Parameters {
		list[i] = param
	}
	return list
}

type Callback struct {
	Name       string
	Parameters []Parameter
	ReturnType Type
	Security   string
	Tags       `json:",omitempty"`
}

func (member *Callback) GetMemberType() string      { return "Callback" }
func (member *Callback) GetName() string            { return member.Name }
func (member *Callback) GetReturnType() rbxapi.Type { return member.ReturnType }
func (member *Callback) GetSecurity() string        { return member.Security }
func (member *Callback) GetParameters() []rbxapi.Parameter {
	list := make([]rbxapi.Parameter, len(member.Parameters))
	for i, param := range member.Parameters {
		list[i] = param
	}
	return list
}

type Parameter struct {
	Type    Type
	Name    string
	Default *string `json:",omitempty"`
}

func (param Parameter) GetType() rbxapi.Type { return param.Type }
func (param Parameter) GetName() string      { return param.Name }
func (param Parameter) GetDefault() (value string, ok bool) {
	if param.Default != nil {
		return *param.Default, true
	}
	return "", false
}

type Enum struct {
	Name  string
	Items []*EnumItem
	Tags  `json:",omitempty"`
}

func (enum *Enum) GetName() string { return enum.Name }
func (enum *Enum) GetItems() []rbxapi.EnumItem {
	list := make([]rbxapi.EnumItem, len(enum.Items))
	for i, item := range enum.Items {
		list[i] = item
	}
	return list
}
func (enum *Enum) GetItem(name string) rbxapi.EnumItem {
	for _, item := range enum.Items {
		if item.GetName() == name {
			return item
		}
	}
	return nil
}

type EnumItem struct {
	Name  string
	Value int
	Tags  `json:",omitempty"`
}

func (item *EnumItem) GetName() string { return item.Name }
func (item *EnumItem) GetValue() int   { return item.Value }

type Tags []string

func (tags Tags) GetTag(tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
func (tags Tags) LenTags() int {
	return len(tags)
}
func (tags *Tags) SetTag(tag ...string) {
	*tags = append(*tags, tag...)
loop:
	for i, n := 1, len(*tags); i < n; {
		for j := 0; j < i; j++ {
			if (*tags)[i] == (*tags)[j] {
				*tags = append((*tags)[:i], (*tags)[i+1:]...)
				n--
				continue loop
			}
		}
		i++
	}
}
func (tags *Tags) UnsetTag(tag ...string) {
loop:
	for i, n := 0, len(*tags); i < n; {
		for j := 0; j < len(tag); j++ {
			if (*tags)[i] == tag[j] {
				*tags = append((*tags)[:i], (*tags)[i+1:]...)
				n--
				continue loop
			}
		}
		i++
	}
}
func (tags Tags) GetTags() []string {
	list := make([]string, 0, len(tags))
	copy(list, tags)
	return list
}

type Type struct {
	Category string
	Name     string
}

func (typ Type) GetName() string     { return typ.Name }
func (typ Type) GetCategory() string { return typ.Category }
func (typ Type) String() string {
	if typ.Category == "" {
		return typ.Name
	}
	return typ.Category + ":" + typ.Name
}
