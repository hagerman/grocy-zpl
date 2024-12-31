^XA

^FX - encoding (28 = UTF-8)
^CI28

^FX - reset position
^LH0,0
^LT0
^LS0

{{ if lt (len .Name) 30 }}

^FX - field for the element 'Line1'
^FO16,44,2
^FWN
^A0,48^FD{{.Name}}^FS

{{ else }}

^FX - field for the element 'Line1'
^FO16,16,2
^FWN
^A0,48^FD{{slice .Name  0 30}}^FS

^FX - field for the element 'Line2'
^FO16,72,2
^FWN
^A0,48^FD{{slice .Name  30 (len .Name)}}^FS

{{ end }}

^FX - field for the element 'Barcode'
^FO24,136,2
^FWN
^BY2,2,64
^BCN,96,N,N^FD{{.Barcode}}^FS

{{ if afterEpoch .DueDate }}

^FX - field for the element 'Expires'
^FO16,248,2
^FWN
^A0,32^FDExp {{.DueDate.Format "2006-01-02"}}^FS

{{ end }}

^XZ