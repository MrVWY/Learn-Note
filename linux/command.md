1、查看某个进程的线程  
eval top -p $(echo $(pidof name)|tr ' ' ',')  
eval top -H -p $(echo $(pidof name)|tr ' ' ',')