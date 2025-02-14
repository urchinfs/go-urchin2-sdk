接口概要说明：

集群所有节点为对等节点，一个接口调用对任意节点触发调用结果均应相同或相似，但因后续内外网IP配置、节点代理的异构存储等不同所以接口调用亦有不同
1. 初始化客户端
NewClient("192.168.242.42:5001")
--- 预埋IP，后续申请外网ip作为网关使用，此节点用于动态获取其他分中心信息，也可从此节点处理网络不通的外网请求

2. 检查是否能连接到ipfs节点：
client.SwarmConnect --- 以节点peer id为参数

3. 上传文件夹：
client.AddDir --- 参考Add接口

4. 上传文件：
client.Add --- 举例：调用add(path)即将数据上传到调用点存储，可能涉及网关存储（本地磁盘）、异构存储obs、minio等，此时集群有该cid数据1份
           --- 举例：调用add(cid, centerId<2>)即将数据上传到调用点存储，并后续同步到center id为2（云脑2）的分中心存储中，此时集群有该cid数据2份

5. 对一个文件或者文件夹离线计算其cid：
cid.GetCid

6. 对一个文件或者文件夹打包生成ipfs car文件：
car.PackCarFormat --- 主要将文件夹打包为一个文件方便后续处理和效率，非必须

7. 将6中生成的car文件导入到ipfs节点：
client.DagImport --- 如果是文件夹则会理论上加速网络上传效率

8. 根据cid从ipfs节点导出car文件：
client.DagExport --- 如果是文件夹则会理论上加速网络上传效率

9. 将8中导出的car文件unpack恢复本身代表的文件或文件夹
car.UnpackCarFormat --- 恢复数据本身形态

10. 根据cid从ipfs节点下载文件或者文件夹：
client.Get

11. 根据cid查询集群中存在该cid数据的所有分中心
CidQuery --- 可用于知晓cid分布情况，避免多次重复上传，即数据已存在则无需再上传（配合5中提前计算出cid）

12. 查询调用节点的peer信息，返回节点id、节点代理存储相关信息等
PeerSelf

13. 根据peer id查询peer信息，返回节点id、节点代理存储相关信息等
PeerQuery

14. 返回集群现时可用的所有节点信息，包括节点id、节点代理存储相关信息等
PeerAll

15. 根据cid查询特定分中心（center id）是否存在该数据
CidExistInCenter

16. 单独调用进行数据同步，参数cid和center id代表将cid数据同步到center id的分中心
CidSync --- 该接口很少单独调用，除非在add、addDir接口未指定同步参数且需要同步的场景，即如果必须优先在add、addDir中一个调用完成上传+同步操作
        --- 举例：在网关代理的节点（本地磁盘）上传，但数据需要存储在云脑2代理的节点，可调用add(path, centerId<2>)即将数据上传到调用点并同步到云脑2分中心
        --- 举例：在中原代理节点（中原obs）上传，后续调用CidSync(cid, centerId<2>)，即同步cid数据到云脑2分中心

17. 异步查询同步操作结果
CheckSyncStatus --- 异步查询，配合add、addDir接口（指定同步参数场景下） 或者 CidSync接口， Failed或者Succeeded为最终状态

18. 触发迁移操作，参数cid、center_id、datastore_path表示将cid数据迁移到center_id分中心存储的urchincache桶的datastore_path路径中
CidMigrate --- 该操作主要用于网络不通而触发迁移，类似现网环境1.0的cache操作，会还原数据到异构存储的datastore_path路径，否则直接Get下载数据即可

19. 异步查询迁移操作结果
CheckMigrateStatus --- Failed或者Succeeded为最终状态

20. 根据cid从peer节点遍历该cid对应的文件或文件夹
List --- 遍历cid下的文件或者文件夹，返回子cid、文件或者文件夹名、文件大小、文件类型
其中文件类型Type中，1为文件夹；2为文件；4为符号链接。

===============
概述：
1. 上传，可以 client.Add上传文件或者 client.AddDir文件夹，也可以对文件夹先car.PackCarFormat打包生成ipfs car文件，再client.DagImport导入

2. 下载，可以直接client.Get文件或文件夹，也可以先对文件夹client.DagExport导出car文件，再car.UnpackCarFormat恢复文件夹

3. 初始化客户端的IP目前预埋网关IP，可以调用接口PeerAll获取现时可用的所有节点信息，后续根据场景和网络可分别请求对应节点

4. 上传前尽可能调用GetCid离线计算cid，这样可用提前查看是否已经存在该cid，对应存储是否有该cid数据等
