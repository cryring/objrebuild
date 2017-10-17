package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Obj is used to rebuild obj file
type Obj struct {
	Mtllib string
	V      []string
	Vn     []string
	Vt     []string
	F      map[string][]string
}

// NewObj is used to create obj
func NewObj() *Obj {
	return &Obj{
		F: make(map[string][]string),
	}
}

// Load is used to read obj file and rebuild
func (o *Obj) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var (
		geoMtl      string
		geoName     string
		fMultiline  bool
		vMultiline  bool
		vnMultiline bool
		vtMultiline bool
	)

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if strings.HasPrefix(line, "mtllib ") {
			o.Mtllib = line
		} else if strings.HasPrefix(line, "v ") || vMultiline {
			o.V = append(o.V, line)
			vMultiline = o.multiLine(line)
		} else if strings.HasPrefix(line, "vn ") || vnMultiline {
			o.Vn = append(o.Vn, line)
			vnMultiline = o.multiLine(line)
		} else if strings.HasPrefix(line, "vt ") || vtMultiline {
			o.Vt = append(o.Vt, line)
			vtMultiline = o.multiLine(line)
		} else if strings.HasPrefix(line, "usemtl ") {
			geoMtl = strings.Replace(line, "usemtl ", "", -1)
		} else if strings.HasPrefix(line, "g ") {
			geoName = strings.Replace(line, "g ", "", -1)
			geoMtl = geoName
		} else if strings.HasPrefix(line, "f ") || fMultiline {
			fkey := geoMtl

			if fkey == "" {
				return fmt.Errorf("illegal obj file")
			}

			if _, ok := o.F[fkey]; !ok {
				o.F[fkey] = make([]string, 0, 32)
			}
			o.F[fkey] = append(o.F[fkey], line)
			fMultiline = o.multiLine(line)
		}
	}
	return nil
}

// Save is used to dump obj to file
func (o *Obj) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(o.Mtllib)
	w.WriteString("\r\n")
	for _, line := range o.V {
		w.WriteString(line)
	}
	w.WriteString("\r\n")
	for _, line := range o.Vn {
		w.WriteString(line)
	}
	w.WriteString("\r\n")
	for _, line := range o.Vt {
		w.WriteString(line)
	}
	w.WriteString("\r\n")
	for mtl, lines := range o.F {
		title := fmt.Sprintf("g %susemtl %s", mtl, mtl)
		w.WriteString(title)
		for _, line := range lines {
			w.WriteString(line)
		}
		w.WriteString("\r\n")
	}
	w.Flush()
	return nil
}

func (o *Obj) multiLine(line string) bool {
	return (strings.HasSuffix(line, "\\\r\n") || strings.HasSuffix(line, "\\\n"))
}
