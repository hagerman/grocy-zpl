^XA

{{ $productName := .Webhook.Name }}
{{ $printProductName := index .ServiceCall.ProductUserFields "product_print_name" }}
{{ if gt (len $printProductName) 0 }}
    {{ $productName = $printProductName }}
{{ end }}
{{ $extras := split (index .ServiceCall.ProductUserFields "product_print_attributes") "," }}

^FX - encoding (28 = UTF-8)
^CI28

^FX - reset position
^LH0,0
^LT0
^LS0

{{ if eq .MediaReady "na_index-2.25x1_2.25x1in" }}

    {{ if lt (len $productName) 30 }}

    ^FX - field for the element 'Line1'
    ^FO16,44,2
    ^FWN
    ^A0,48^FD{{$productName}}^FS

    {{ else }}

    ^FX - field for the element 'Line1'
    ^FO16,16,2
    ^FWN
    ^A0,48^FD{{slice $productName  0 30}}^FS

    ^FX - field for the element 'Line2'
    ^FO16,72,2
    ^FWN
    ^A0,48^FD{{slice $productName  30 (len $productName)}}^FS

    {{ end }}

    ^FX - field for the element 'Barcode'
    ^FO24,136,2
    ^FWN
    ^BY2,2,64
    ^BCN,96,N,N^FD{{.Webhook.Barcode}}^FS

    {{ if afterEpoch .Webhook.DueDate }}

    ^FX - field for the element 'Expires'
    ^FO16,248,2
    ^FWN
    ^A0,32^FDExp {{.Webhook.DueDate.Format "2006-01-02"}}^FS

    {{ end }}

    ^FX - field for the element 'Product Group Print Name'
    ^FO651,248,1
    ^FWN
    ^A0,32^FD{{index .ServiceCall.ProductGroupUserFields "product_group_print_name"}}^FS

{{ else if eq .MediaReady "na_index-4x6_4x6in" }}

    ^FWR

    ^CF0,100
    ^FO450,150
    ^FB800,4^FD{{uppercase $productName}}^FS

    {{ if gt (len $extras) 0 }}

        ^CFE,50
        ^FO250,150^FB800,1^FD{{ index $extras 0 }}{{ if gt (len $extras) 0 }} / {{ index $extras 1 }} {{ end }}^FS
        ^FO225,150
        ^GB0,700,3,B^FS

    {{ end }}

    {{ if gt (len $extras) 2 }}

        ^CFE,50
        ^FO350,150^FB800,1^FD{{ index $extras 2 }}{{ if gt (len $extras) 3 }} / {{ index $extras 3 }} {{ end }}^FS
        ^FO325,150
        ^GB0,700,3,B^FS

    {{ end }}

    ^CF0,50
    ^FO150,150^FD{{index .ServiceCall.ProductGroupUserFields "product_group_print_name"}}^FS

    ^FWI
    {{ if lt (len .Webhook.Name) 30 }}
        ^CF0,60
        ^FB600,1
        ^FO460,1130^FD{{uppercase $productName}}^FS
    {{ else }}
        ^CF0,60
        ^FB550,2
        ^FO510,1100^FD{{uppercase $productName}}^FS
    {{ end }}

    {{ if afterEpoch .Webhook.DueDate }}
        ^CF0,25
        ^FB600,1,0,R
        ^FO160,1170^FDEXPIRES^FS
        ^FB600,1,0,R
        ^CF0,40
        ^FO160,1130^FD{{.Webhook.DueDate.Format "2006-01-02"}}^FS
        ^FWN
    {{ end }}

    ^BY2,2,64^FO300,1540^BCN,96,N^FD{{.Webhook.Barcode}}^FS

{{ end }}

^XZ