package pak

type FileEntry struct {
	Offset  uint32
	Length  uint32
	Data    []byte `struct:"-"`
	Name    string `struct:"-"`
	ID      int    `struct:"-"`
	Replace bool   `struct:"-"`
}

func (e *FileEntry) OpenScript() {

}
