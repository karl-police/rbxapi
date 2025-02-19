package rbxapidump

import (
	"github.com/karl-police/rbxapi"
	"github.com/karl-police/rbxapi/patch"
)

// copyClass returns a deep copy of a generic rbxapi.Class.
func copyClass(class rbxapi.Class) *Class {
	members := class.GetMembers()
	c := Class{
		Name:       class.GetName(),
		Superclass: class.GetSuperclass(),
		Members:    make([]rbxapi.Member, 0, len(members)),
		Tags:       Tags(class.GetTags()),
	}
	for _, member := range members {
		if member := copyMember(member); member != nil {
			switch member := member.(type) {
			case *Property:
				member.Class = class.GetName()
			case *Function:
				member.Class = class.GetName()
			case *Event:
				member.Class = class.GetName()
			case *Callback:
				member.Class = class.GetName()
			}
			c.Members = append(c.Members, member)
		}
	}
	return &c
}

// copyMember returns a deep copy of a generic rbxapi.Member.
func copyMember(member rbxapi.Member) rbxapi.Member {
	switch member := member.(type) {
	case rbxapi.Property:
		return &Property{
			Name:      member.GetName(),
			ValueType: copyType(member.GetValueType()),
			Tags:      Tags(member.GetTags()),
		}
	case rbxapi.Function:
		// Function and Callback have the same methods.
		switch member.GetMemberType() {
		case "Function":
			return &Function{
				Name:       member.GetName(),
				ReturnType: copyType(member.GetReturnType()),
				Parameters: copyParameters(member.GetParameters()),
				Tags:       Tags(member.GetTags()),
			}
		case "Callback":
			return &Callback{
				Name:       member.GetName(),
				ReturnType: copyType(member.GetReturnType()),
				Parameters: copyParameters(member.GetParameters()),
				Tags:       Tags(member.GetTags()),
			}
		}
	case rbxapi.Event:
		return &Event{
			Name:       member.GetName(),
			Parameters: copyParameters(member.GetParameters()),
			Tags:       Tags(member.GetTags()),
		}
	}
	return nil
}

// copyEnum returns a deep copy of a generic rbxapi.Enum.
func copyEnum(enum rbxapi.Enum) *Enum {
	items := enum.GetEnumItems()
	e := Enum{
		Name:  enum.GetName(),
		Items: make([]*EnumItem, 0, len(items)),
		Tags:  Tags(enum.GetTags()),
	}
	for _, item := range items {
		if item := copyEnumItem(item); item != nil {
			item.Enum = enum.GetName()
			e.Items = append(e.Items, item)
		}
	}
	return &e
}

// copyEnumItem returns a deep copy of a generic rbxapi.EnumItem.
func copyEnumItem(item rbxapi.EnumItem) *EnumItem {
	return &EnumItem{
		Name:  item.GetName(),
		Value: item.GetValue(),
		Tags:  item.GetTags(),
	}
}

// copyParameters returns a deep copy of a list of generic rbxapi.Parameter
// values.
func copyParameters(params rbxapi.Parameters) []Parameter {
	list := make([]Parameter, params.GetLength())
	for i := 0; i < len(list); i++ {
		param := params.GetParameter(i)
		list[i].Type.SetFromType(param.GetType())
		list[i].Name = param.GetName()
		list[i].Default, list[i].HasDefault = param.GetDefault()
	}
	return list
}

// copyType returns a deep copy of a generic rbxapi.Type.
func copyType(typ rbxapi.Type) Type {
	var t Type
	t.SetFromType(typ)
	return t
}

// Patch transforms the API structure by applying a list of patch actions.
//
// Patch implements the patch.Patcher interface.
func (root *Root) Patch(actions []patch.Action) {
	for i, action := range actions {
		if action, ok := action.(patch.Member); ok {
			if aclass, amember := (action.GetClass()), (action.GetMember()); aclass != nil && amember != nil {
				name := aclass.GetName()
				for _, class := range root.Classes {
					if class.Name == name {
						class.Patch(actions[i : i+1])
						break
					}
				}
				continue
			}
		}
		if action, ok := action.(patch.Class); ok {
			if aclass := action.GetClass(); aclass != nil {
				switch action.GetType() {
				case patch.Remove:
					name := aclass.GetName()
					for i, class := range root.Classes {
						if class.Name == name {
							copy(root.Classes[i:], root.Classes[i+1:])
							root.Classes[len(root.Classes)-1] = nil
							root.Classes = root.Classes[:len(root.Classes)-1]
							break
						}
					}
				case patch.Add:
					root.Classes = append(root.Classes, copyClass(aclass))
				case patch.Change:
					name := aclass.GetName()
					for _, class := range root.Classes {
						if class.Name == name {
							class.Patch(actions[i : i+1])
							break
						}
					}
				}
				continue
			}
		}
		if action, ok := action.(patch.EnumItem); ok {
			if aenum, aitem := (action.GetEnum()), (action.GetEnumItem()); aenum != nil && aitem != nil {
				name := aenum.GetName()
				for _, enum := range root.Enums {
					if enum.Name == name {
						enum.Patch(actions[i : i+1])
						break
					}
				}
				continue
			}
		}
		if action, ok := action.(patch.Enum); ok {
			if aenum := action.GetEnum(); aenum != nil {
				switch action.GetType() {
				case patch.Remove:
					name := aenum.GetName()
					for i, enum := range root.Enums {
						if enum.Name == name {
							copy(root.Enums[i:], root.Enums[i+1:])
							root.Enums[len(root.Enums)-1] = nil
							root.Enums = root.Enums[:len(root.Enums)-1]
							break
						}
					}
				case patch.Add:
					root.Enums = append(root.Enums, copyEnum(aenum))
				case patch.Change:
					name := aenum.GetName()
					for _, enum := range root.Enums {
						if enum.Name == name {
							enum.Patch(actions[i : i+1])
							break
						}
					}
				}
				continue
			}
		}
	}
}

func (class *Class) Patch(actions []patch.Action) {
	for i, action := range actions {
		if action, ok := action.(patch.Member); ok {
			if aclass, amember := (action.GetClass()), (action.GetMember()); aclass != nil && amember != nil {
				switch action.GetType() {
				case patch.Remove:
					name := amember.GetName()
					for i, member := range class.Members {
						if member.GetName() == name {
							copy(class.Members[i:], class.Members[i+1:])
							class.Members[len(class.Members)-1] = nil
							class.Members = class.Members[:len(class.Members)-1]
							break
						}
					}
				case patch.Add:
					if member := copyMember(amember); member != nil {
						class.Members = append(class.Members, member)
					}
				case patch.Change:
					name := amember.GetName()
					mtype := amember.GetMemberType()
					for _, member := range class.Members {
						if member.GetName() == name && member.GetMemberType() == mtype {
							if member, ok := member.(patch.Patcher); ok {
								member.Patch(actions[i : i+1])
							}
							break
						}
					}
				}
				continue
			}
		}
		if action, ok := action.(patch.Class); ok {
			if aclass := action.GetClass(); aclass != nil {
				if action.GetType() != patch.Change {
					continue
				}
				switch action.GetField() {
				case "Name":
					if v, ok := action.GetNext().(string); ok {
						class.Name = v
					}
				case "Superclass":
					if v, ok := action.GetNext().(string); ok {
						class.Superclass = v
					}
				case "Tags":
					if v, ok := action.GetNext().([]string); ok {
						class.Tags = Tags(Tags(v).GetTags())
					}
				}
				continue
			}
		}
	}
}

func (member *Property) Patch(actions []patch.Action) {
	for _, action := range actions {
		if action.GetType() != patch.Change {
			continue
		}
		switch action.GetField() {
		case "Name":
			if v, ok := action.GetNext().(string); ok {
				member.Name = v
			}
		case "ValueType":
			switch v := action.GetNext().(type) {
			case rbxapi.Type:
				member.ValueType.SetFromType(v)
			case string:
				member.ValueType = Type(v)
			}
		case "Tags":
			if v, ok := action.GetNext().([]string); ok {
				member.Tags = Tags(Tags(v).GetTags())
			}
		}
	}
}

func (member *Function) Patch(actions []patch.Action) {
	for _, action := range actions {
		if action.GetType() != patch.Change {
			continue
		}
		switch action.GetField() {
		case "Name":
			if v, ok := action.GetNext().(string); ok {
				member.Name = v
			}
		case "ReturnType":
			switch v := action.GetNext().(type) {
			case rbxapi.Type:
				member.ReturnType.SetFromType(v)
			case string:
				member.ReturnType = Type(v)
			}
		case "Parameters":
			if v, ok := action.GetNext().(rbxapi.Parameters); ok {
				member.Parameters = copyParameters(v)
			}
		case "Tags":
			if v, ok := action.GetNext().([]string); ok {
				member.Tags = Tags(Tags(v).GetTags())
			}
		}
	}
}

func (member *Event) Patch(actions []patch.Action) {
	for _, action := range actions {
		if action.GetType() != patch.Change {
			continue
		}
		switch action.GetField() {
		case "Name":
			if v, ok := action.GetNext().(string); ok {
				member.Name = v
			}
		case "Parameters":
			if v, ok := action.GetNext().(rbxapi.Parameters); ok {
				member.Parameters = copyParameters(v)
			}
		case "Tags":
			if v, ok := action.GetNext().([]string); ok {
				member.Tags = Tags(Tags(v).GetTags())
			}
		}
	}
}

func (member *Callback) Patch(actions []patch.Action) {
	for _, action := range actions {
		if action.GetType() != patch.Change {
			continue
		}
		switch action.GetField() {
		case "Name":
			if v, ok := action.GetNext().(string); ok {
				member.Name = v
			}
		case "ReturnType":
			switch v := action.GetNext().(type) {
			case rbxapi.Type:
				member.ReturnType.SetFromType(v)
			case string:
				member.ReturnType = Type(v)
			}
		case "Parameters":
			if v, ok := action.GetNext().(rbxapi.Parameters); ok {
				member.Parameters = copyParameters(v)
			}
		case "Tags":
			if v, ok := action.GetNext().([]string); ok {
				member.Tags = Tags(Tags(v).GetTags())
			}
		}
	}
}

func (enum *Enum) Patch(actions []patch.Action) {
	for i, action := range actions {
		if action, ok := action.(patch.EnumItem); ok {
			if aenum, aitem := (action.GetEnum()), (action.GetEnumItem()); aenum != nil && aitem != nil {
				switch action.GetType() {
				case patch.Remove:
					name := aitem.GetName()
					for i, item := range enum.Items {
						if item.GetName() == name {
							copy(enum.Items[i:], enum.Items[i+1:])
							enum.Items[len(enum.Items)-1] = nil
							enum.Items = enum.Items[:len(enum.Items)-1]
							break
						}
					}
				case patch.Add:
					enum.Items = append(enum.Items, copyEnumItem(aitem))
				case patch.Change:
					name := aitem.GetName()
					for _, item := range enum.Items {
						if item.GetName() == name {
							item.Patch(actions[i : i+1])
							break
						}
					}
				}
				continue
			}
		}
		if action, ok := action.(patch.Enum); ok {
			if aenum := action.GetEnum(); aenum != nil {
				if action.GetType() != patch.Change {
					continue
				}
				switch action.GetField() {
				case "Name":
					if v, ok := action.GetNext().(string); ok {
						enum.Name = v
					}
				case "Tags":
					if v, ok := action.GetNext().([]string); ok {
						enum.Tags = Tags(Tags(v).GetTags())
					}
				}
				continue
			}
		}
	}
}

func (item *EnumItem) Patch(actions []patch.Action) {
	for _, action := range actions {
		if action.GetType() != patch.Change {
			continue
		}
		switch action.GetField() {
		case "Name":
			if v, ok := action.GetNext().(string); ok {
				item.Name = v
			}
		case "Value":
			if v, ok := action.GetNext().(int); ok {
				item.Value = v
			}
		case "Tags":
			if v, ok := action.GetNext().([]string); ok {
				item.Tags = Tags(Tags(v).GetTags())
			}
		}
	}
}
