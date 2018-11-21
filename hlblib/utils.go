package hlblib

func PanicIfErrorExist(err error) {
	if err != nil {
		panic(err)
	}
}
