^XA

^LH0,0
^LT0
^LS0

{{ if lt (len .Name) 30 }}

^FXfield for the element 'Line1'
^FO16,44,2
^FWN
^A0,48^FD{{.Name}}^FS

{{ else }}

^FXfield for the element 'Line1'
^FO16,16,2
^FWN
^A0,48^FD{{slice .Name  0 30}}^FS

^FXfield for the element 'Line2'
^FO16,72,2
^FWN
^A0,48^FD{{slice .Name  30 (len .Name)}}^FS

{{ end }}

^FXfield for the element 'Barcode'
^FO24,136,2
^FWN
^BY2,2,64
^BCN,96,N,N^FD{{.Barcode}}^FS

^FXfield for the element 'Expires'
^FO16,248,2
^FWN
^A0,32^FDExp {{.DueDate.Format "2006-01-02"}}^FS
^XZ