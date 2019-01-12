@echo off
echo 清除所有obj pch idb pdb ncb opt plg res sbr ilk suo文件，请稍等......
pause
del /f /s /q .\*.obj
del /f /s /q .\*.pch
del /f /s /q .\*.idb
del /f /s /q .\*.pdb
del /f /s /q .\*.ncb 
del /f /s /q .\*.opt 
del /f /s /q .\*.plg
del /f /s /q .\*.sdf
del /f /s /q .\*.sbr
del /f /s /q .\*.ilk
del /f /s /q .\*.aps
del /f /s /q .\*.ipch
del /f /s /q .\*.dmp
del /f /s /q .\*.log
del /f /s /q .\*.err
del /f /s /q .\*.DS_Store
del /f /s /q server.exe
del /f /s /q client.exe
del /f /s /q nohup.out
del /f /s /q server
del /f /s /q client
rd  /s /q .\pkg
rd  /s /q .\.idea



echo 清除文件完成！
echo. & pause