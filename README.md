# grocy-zpl
A way to print labels from Grocy to Zebra printers. This is known to support the Zebra ZSB DP24.

## Instructions
Set the `PRINTER_URL` environment variable to be the IPP URL of your printer.
For example `http://192.168.1.228:631/ipp/print`.
There is an included template for 2.25x1.00 in labels, but this can be overriden
if needed by setting the `TEMPLATE_PATH` variable.