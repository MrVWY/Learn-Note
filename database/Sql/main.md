## group by、having、order by

    from-->where-->group by -->having --> select--- >order by;

使用顺序：group by 、having、order by(排序)

group by : 按照某个字段或者某些字段进行分组

having : 对分组之后的数据进行再次过滤

having是在分好组后找出特定的分组，通常是以筛选聚合函数的结果

聚合函数：

1. MIN 最小值 
2. MAX 最大值 
3. SUM 求和 
4. AVG 求平均 
5. COUNT 计数


### where与having的根本区别在于：

where子句在Group By分组和聚合函数之前对数据行进行过滤

having子句对Group BY分组和聚合函数之后的数据行进行过滤。
