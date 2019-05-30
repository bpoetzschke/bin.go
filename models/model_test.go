package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordList_Diff(t *testing.T) {
	a := WordList{
		{
			Value: "1",
		},
		{
			Value: "2",
		},
		{
			Value: "3",
		},
	}

	b := WordList{
		{
			Value: "2",
		},
	}

	expected := WordList{
		a[0], a[2],
	}

	res := a.Diff(b)

	require.EqualValues(t, expected, res)

	res_2 := b.Diff(a)
	require.EqualValues(t, expected, res_2)
}
