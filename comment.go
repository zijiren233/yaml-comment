package yamlcomment

import (
	"reflect"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const (
	HeadCommentTag = "hc="
	LineCommentTag = "lc="
	FootCommentTag = "fc="
	OmitemptyTag   = "omitempty"
	InlineTag      = "inline"
	FlowTag        = "flow"
)

type option struct {
	fieldName string
	omitempty bool
	skip      bool
	inline    bool
	flow      bool
}

type comment struct {
	HeadComment string
	LineComment string
	FootComment string
}

func setComment(cm *comment, parts ...string) {
	if cm == nil {
		return
	}
	var pre *string
	for _, part := range parts {
		if strings.HasPrefix(part, HeadCommentTag) {
			cm.HeadComment = strings.TrimPrefix(part, HeadCommentTag)
			pre = &cm.HeadComment
		} else if strings.HasPrefix(part, LineCommentTag) {
			cm.LineComment = strings.TrimPrefix(part, LineCommentTag)
			pre = &cm.LineComment
		} else if strings.HasPrefix(part, FootCommentTag) {
			cm.FootComment = strings.TrimPrefix(part, FootCommentTag)
			pre = &cm.FootComment
		} else if pre != nil {
			*pre += "," + part
		}
	}
}

func newComment(parts ...string) *comment {
	cm := new(comment)
	setComment(cm, parts...)
	return cm
}

type CommentEncoder struct {
	encoder *yaml.Encoder
}

func NewEncoder(encoder *yaml.Encoder) *CommentEncoder {
	return &CommentEncoder{
		encoder: encoder,
	}
}

func (e *CommentEncoder) Encode(v any) error {
	node, err := AnyToYamlNode(v)
	if err != nil {
		return err
	}
	return e.encoder.Encode(node)
}

func Marshal(v any) ([]byte, error) {
	node, err := AnyToYamlNode(v)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(node)
}

func isZero(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	return value.IsZero()
}

func isNil(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func parseTags(tag string) (*option, *comment) {
	parts := strings.Split(tag, ",")

	var op = &option{
		fieldName: parts[0],
	}

	if op.fieldName == "-" {
		op.skip = true
		return op, nil
	}

	parts = parts[1:]

	comments := make([]string, 0, len(parts))

	for _, part := range parts {
		switch part {
		case OmitemptyTag:
			op.omitempty = true
		case InlineTag:
			op.inline = true
		case FlowTag:
			op.flow = true
		default:
			comments = append(comments, part)
		}
	}

	return op, newComment(comments...)
}

func AnyToYamlNode(model any) (*yaml.Node, error) {
	if n, ok := model.(*yaml.Node); ok {
		return n, nil
	}

	if m, ok := model.(yaml.Marshaler); ok && !isNil(reflect.ValueOf(model)) {
		res, err := m.MarshalYAML()
		if err != nil {
			return nil, err
		}

		if n, ok := res.(*yaml.Node); ok {
			return n, nil
		}

		model = res
	}

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	node := new(yaml.Node)

	switch v.Kind() {
	case reflect.Struct:
		node.Kind = yaml.MappingNode

		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanInterface() {
				continue
			}

			tag := t.Field(i).Tag.Get("yaml")

			op, cm := parseTags(tag)

			if op.skip || (op.omitempty && isZero(field)) {
				continue
			}

			if op.fieldName == "" {
				op.fieldName = strings.ToLower(t.Field(i).Name)
			}

			var value any
			if field.CanInterface() {
				value = field.Interface()
			}

			var style yaml.Style
			if op.flow {
				style |= yaml.FlowStyle
			}

			if op.inline {
				child, err := AnyToYamlNode(value)
				if err != nil {
					return nil, err
				}

				if child.Kind == yaml.MappingNode || child.Kind == yaml.SequenceNode {
					appendNodes(node, child.Content...)
				}
			} else if err := addToMap(node, op.fieldName, value, cm, style); err != nil {
				return nil, err
			}
		}
	case reflect.Map:
		node.Kind = yaml.MappingNode
		keys := v.MapKeys()
		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			if err := addToMap(node, k.Interface(), v.MapIndex(k).Interface(), nil, 0); err != nil {
				return nil, err
			}
		}
	case reflect.Slice:
		node.Kind = yaml.SequenceNode
		nodes := make([]*yaml.Node, v.Len())

		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)

			var err error

			nodes[i], err = AnyToYamlNode(element.Interface())
			if err != nil {
				return nil, err
			}
		}
		appendNodes(node, nodes...)
	default:
		if err := node.Encode(model); err != nil {
			return nil, err
		}
	}

	return node, nil
}

func appendNodes(dest *yaml.Node, nodes ...*yaml.Node) {
	if dest.Content == nil {
		dest.Content = nodes
		return
	}

	dest.Content = append(dest.Content, nodes...)
}

func addToMap(dest *yaml.Node, fieldName, in any, cm *comment, style yaml.Style) error {
	key, err := AnyToYamlNode(fieldName)
	if err != nil {
		return err
	}

	value, err := AnyToYamlNode(in)
	if err != nil {
		return err
	}
	value.Style = style

	addComment(key, cm)
	appendNodes(dest, key, value)

	return nil
}

func addComment(node *yaml.Node, cm *comment) {
	if cm == nil {
		return
	}

	node.HeadComment = cm.HeadComment
	node.LineComment = cm.LineComment
	node.FootComment = cm.FootComment
}
