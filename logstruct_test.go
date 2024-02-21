package logstruct

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	r := require.New(t)

	type EmailBox struct {
		Email       string `log:"email"`
		HiddenEmail string `log:"-"`
	}

	type Data struct {
		Contract    int     `db:"c_id" log:"c_id"`
		OperLogin   string  `json:"oper_login" log:"oper_login"`
		Balance     float64 `log:"balance"`
		Exists      bool    `log:"exists"`
		NotExists1  bool    `log:"-"` // пропускается
		NotExists2  bool    // тоже самое что и выше, если нет log тэга - не попадет в logrus.Fields
		TestPointer *int    `log:"test_pointer"` // поинтер значения или структуры пропускаются из лога
		EmailBox
		TestDate         time.Time         `log:"test_date"`
		SliceInt         []int             `log:"slice_int"`
		SliceStr         []string          `log:"slice_str"`
		MapStr           map[string]string `log:"map_str"`
		EmptySlice       []int             `log:"empty_slice"`
		NullString       NullString        `log:"null_string"`
		NullTime         NullTime          `log:"null_time"`
		NullFloat64      NullFloat64       `log:"null_float64"`
		NullInt64        NullInt64         `log:"null_int64"`
		NullBool         NullBool          `log:"null_bool"`
		EmptyNullString  NullString        `log:"empty_null_string"`
		EmptyNullTime    NullTime          `log:"empty_null_time"`
		EmptyNullFloat64 NullFloat64       `log:"empty_null_float64"`
		EmptyNullInt64   NullInt64         `log:"empty_null_int64"`
		EmptyNullBool    NullBool          `log:"empty_null_bool"`
	}

	i := 10
	testDate := time.Date(2023, 12, 26, 0, 0, 0, 0, time.Local)

	d := Data{
		Contract:    1234,
		OperLogin:   "m.zarif@sarkor.uz",
		Balance:     12345.6,
		Exists:      true,
		NotExists1:  true,
		NotExists2:  true,
		TestPointer: &i,
		EmailBox: EmailBox{
			Email:       "test@email.uz",
			HiddenEmail: "hiddentest@email.uz",
		},
		TestDate: testDate,
		SliceInt: []int{1, 2, 3},
		SliceStr: []string{"a", "b", "c"},
		MapStr: map[string]string{
			"1": "test1",
			"2": "test2",
		},
		NullString:  NewString("null_string"),
		NullTime:    NewNullTime(testDate),
		NullFloat64: NewFloat64(1),
		NullInt64:   NewInt64(1),
		NullBool:    NullBool{Bool: true, Valid: true},
	}

	expected := logrus.Fields{
		"c_id":               1234,
		"oper_login":         "m.zarif@sarkor.uz",
		"balance":            12345.6,
		"exists":             true,
		"email":              "test@email.uz",
		"test_date":          testDate.String(),
		"slice_int":          "1,2,3",
		"slice_str":          "a,b,c",
		"map_str":            `{"1":"test1","2":"test2"}`,
		"empty_slice":        "",
		"null_string":        "null_string",
		"null_time":          testDate.String(),
		"null_float64":       1.0,
		"null_int64":         int64(1),
		"null_bool":          true,
		"empty_null_string":  "-",
		"empty_null_time":    "-",
		"empty_null_float64": "-",
		"empty_null_int64":   "-",
		"empty_null_bool":    "-",
	}

	// умеет работать как со значением
	result := Fields(d)
	r.Equal(expected, result)

	// // умеет работать так и по поинтеру структуры
	// result = Fields(&d)
	// r.Equal(expected, result)

	// // при передаче не структуры = результат будет пустым
	// result = Fields("test")
	// r.Empty(result)

}
