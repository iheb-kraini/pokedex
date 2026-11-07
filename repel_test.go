package main
import "testing" 
func TestCleanInput(t *testing.T) {
    cases := []struct {
	input    string
	expected []string
}{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	{
		input:    "Charmander Bulbasaur PIKACHU",
		expected: []string{"charmander", "bulbasaur", "pikachu"},
	},

	{
		input:    "a  b",
		expected: []string{"a", "b"},
	},

	{
		input:    "   ",
		expected: []string{},
	},

}






for _, c := range cases {
	actual := cleanInput(c.input)
	// Check the length of the actual slice
	// if they don't match, use t.Errorf and continue to the next case
	if len(actual) != len(c.expected) {
		t.Errorf("length mismatch in case %s",c);
	}
	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
	if (word != expectedWord){
		t.Errorf("expected %s got %s" ,expectedWord,actual);
	}
	}
}









}
