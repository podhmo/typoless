package typoless

import (
	"fmt"
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
	}

	for i, c := range cases {
		c := c
		t.Run(fmt.Sprintf("case:%d", i), func(t *testing.T) {
			got := c.input.TableName()
			if c.want != got {
				t.Errorf("\nwant\n\t%v\nbut\n\t%s", c.want, got)
			}
		})
	}
}
