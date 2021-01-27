echo off & color 0A

rem 参考文章 https://github.com/google/protobuf/blob/master/cmake/README.md
rem 默认当前操作系统已安装 git 和 cmake,并配置好了环境变量

set "WORK_DIR=%cd%"
echo %WORK_DIR%

rem 设置所需要的Protobuf版本,最新版本可以在github上查到 https://github.com/google/protobuf
set "PROTOBUF_VESION=v3.5.0"
echo %PROTOBUF_VESION%
set "PROTOBUF_PATH=protobuf_%PROTOBUF_VESION%"
echo %PROTOBUF_PATH%

rem 从githug上拉取protobuf源代码
git clone -b %PROTOBUF_VESION% https://github.com/google/protobuf.git %PROTOBUF_PATH%

rem 从github上拉取gmock
cd %PROTOBUF_PATH%
git clone -b release-1.7.0 https://github.com/google/googlemock.git gmock

rem 从github上拉取gtest
cd gmock
git clone -b release-1.7.0 https://github.com/google/googletest.git gtest

cd %WORK_DIR%
rem 设置VS工具集,相当于指定VS版本,取决于VS的安装路径
set VS_DEV_CMD="D:\Program Files (x86)\Microsoft Visual Studio 14.0\Common7\Tools\VsDevCmd.bat"
rem 设置工程文件夹名字,用来区分不同的VS版本
set "BUILD_PATH=protobuf_%PROTOBUF_VESION%_vs2015_sln"
echo %BUILD_PATH%
if not exist %BUILD_PATH% md %BUILD_PATH%
cd %BUILD_PATH%
rem 设置编译版本 Debug Or Release
set "MODE=Release"
echo %MODE%
if not exist %MODE% md %MODE%
cd %MODE%
echo %cd%

set "CMAKELISTS_DIR=%WORK_DIR%\%PROTOBUF_PATH%\cmake"
echo %CMAKELISTS_DIR%

rem 开始构建和编译
call %VS_DEV_CMD%
cmake %CMAKELISTS_DIR% -G "NMake Makefiles" -DCMAKE_BUILD_TYPE=%MODE%
call extract_includes.bat
nmake /f Makefile

echo %cd%
echo %PROTOBUF_VESION%
echo %BUILD_PATH%
echo %MODE%
pause