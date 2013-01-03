 set timefmt "%s"
 set xdata time
 set format x "%d/%m\n%H:%M"
 set grid
 f(x)=(x-3600*7)
 #plot "temp.dat" using (f(column("epoch"))):"Temperature" smooth csplines title "xbee"
 #plot "temp.dat" using (f(column("epoch"))):"Temperature" with lines title "xbee"
 plot "temp.dat" using (f(column("epoch"))):"Temperature" title "xbee"
pause 30


