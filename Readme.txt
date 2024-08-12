
1. 初始化客户端
NewClient("192.168.242.42:5001")

2. 检查是否能连接到ipfs节点：
client.SwarmConnect

3. 上传文件夹：
client.AddDir

4. 上传文件：
client.Add

5. 对一个文件或者文件夹离线计算其cid：
cid.GetCid

6. 对一个文件或者文件夹打包生成ipfs car文件：
car.PackCarFormat

7. 将6中生成的car文件导入到ipfs节点：
client.DagImport

8. 根据cid从ipfs节点导出car文件：
client.DagExport

9. 将8中导出的car文件unpack恢复本身代表的文件或文件夹
car.UnpackCarFormat

10. 根据cid从ipfs节点下载文件或者文件夹：
client.Get

===============
概述：
1. 上传，可以 client.Add上传文件或者 client.AddDir文件夹，也可以对文件夹先car.PackCarFormat打包生成ipfs car文件，再client.DagImport导入

2. 下载，可以直接client.Get文件或文件夹，也可以先对文件夹client.DagExport导出car文件，再car.UnpackCarFormat恢复文件夹

3. 初始化客户端的IP计划后续新增接口从服务下发，目前可以先指定固定几个ipfs节点ip



