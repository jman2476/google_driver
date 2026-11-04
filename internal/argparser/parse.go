package argparser

type FlagArg struct {
	flag  string
	value string
}

type Arguments struct {
	Command    string
	SourcePath string
	TargetPath string
	Flags      []FlagArg
}

func ParseFlags(args []string) Arguments {

	return Arguments{}
}
