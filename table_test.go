package typoless

import (
	"fmt"
	"strings"
	"testing"
)

func TestTableName(t *testing.T) {
	cases := []struct {
		input tablelike
		want  string
	}{
		{
			input: Table("xs"),
			want:  "xs",
		},
		{
			input: func() tablelike {
				xs := Table("xs")
				ys := Table("ys")
				return xs.Join(ys, "ON xs.id = ys.x_id")
			}(),
			want: "xs JOIN ys ON xs.id = ys.x_id",
		},
		{
			input: func() tablelike {
				xs := Table("xs")
				ys := Table("ys")

				lhs := xs.As("lhs")
				rhs := ys.As("rhs")
				return lhs.Join(rhs, "ON lhs.id = rhs.x_id")
			}(),
			want: "xs as lhs JOIN ys as rhs ON lhs.id = rhs.x_id",
		},
		{
			input: func() tablelike {
				xTable := struct {
					Table
					ID Int64Field
				}{
					Table: Table("xs"),
					ID:    Int64Field("id"),
				}
				yTable := struct {
					Table
					ID  Int64Field
					XID Int64Field
				}{
					Table: Table("ys"),
					ID:    Int64Field("id"),
					XID:   Int64Field("x_id"),
				}

				lhs := xTable // copy
				Alias(&lhs, &xTable, "lhs")
				rhs := yTable // copy
				Alias(&rhs, &yTable, "rhs")
				return lhs.Join(rhs, On(lhs.ID, rhs.XID))
			}(),
			want: "xs as lhs JOIN ys as rhs ON lhs.id=rhs.x_id",
		},
		{
			input: func() tablelike {
				people := struct {
					Table
					ID       Int64Field
					FatherID Int64Field
					MotherID Int64Field
				}{
					Table:    Table("people"),
					ID:       Int64Field("id"),
					FatherID: Int64Field("father_id"),
					MotherID: Int64Field("mother_id"),
				}

				p := people // copy
				Alias(&p, &people, "p")
				father := people // copy
				Alias(&father, &people, "father")
				mother := people // copy
				Alias(&mother, &people, "mother")

				return p.
					Join(father, On(p.FatherID, father.ID)).
					Join(mother, On(p.MotherID, mother.ID))
			}(),
			want: `
people as p
 JOIN people as father ON p.father_id=father.id
 JOIN people as mother ON p.mother_id=mother.id
`,
		},
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			got := strings.ReplaceAll(c.input.TableName(), "\n", "")
			want := strings.ReplaceAll(c.want, "\n", "")
			if want != got {
				t.Errorf("\nwant\n\t%v\nbut\n\t%s", want, got)
			}
		})
	}
}
