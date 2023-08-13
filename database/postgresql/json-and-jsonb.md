### Json & jsonb 的异同

​	PostgreSQL 支持两种 JSON 数据类型：`json`与`jsonb`。二者在使用上几乎无差异，主要区别是 json 存储的是文本格式，而 jsonb 存储的是二进制格式。因此：

- json 在插入时不需要额外处理，而 jsonb 需要处理为二进制，所以 json 的插入比 jsonb 要快；
- jsonb 以二进制来存储已经解析好的数据，在检索的时候不需要再额外处理，因此检索的性能比 json 要好；
- 另外 jsonb 支持索引，若无特殊需求，推荐使用 jsonb。

​	所以json和jsonb类型在使用上几乎完全一致，两者的区别主要在存储上，json数据类型直接存储输入文本的完全的拷贝，jsonb数据类型以二进制格式进行存储。同时jsonb相较于json更高效，处理速度提升非常大，且**支持索引**

### json and jsonb 操作符

| 操作符 | 右操作数类型 | 结果类型          | 描述                                                         | 例子                                                      | 结果           |
| ------ | ------------ | ----------------- | ------------------------------------------------------------ | --------------------------------------------------------- | -------------- |
| `->`   | `int`        | `json` or `jsonb` | Get JSON array element (indexed from zero, negative integers count from the end)（获得JSON数组元素（索引从零开始）。） | select `'[{"a":"foo"},{"b":"bar"},{"c":"baz"}]'::json->2` | `{"c":"baz"}`  |
| `->`   | `text`       | `json` or `jsonb` | Get JSON object field by key（根据键获得JSON对象的域。）     | select `'{"a": {"b":"foo"}}'::json->'a'`                  | `{"b":"foo"}`  |
| `->>`  | `int`        | `text`            | Get JSON array element as `text`（获得JSON数组元素的文本形式） | select `'[1,2,3]'::json->>2`                              | `3`            |
| `->>`  | `text`       | `text`            | Get JSON object field as `text`（获得JSON对象域的文本形式）  | select `'{"a":1,"b":2}'::json->>'b'`                      | `2`            |
| `#>`   | `text[]`     | `json` or `jsonb` | Get JSON object at the specified path（获得在指定路径上的JSON对象） | select `'{"a": {"b":{"c": "foo"}}}'::json#>'{a,b}'`       | `{"c": "foo"}` |
| `#>>`  | `text[]`     | `text`            | Get JSON object at the specified path as `text`（获得在指定路径上的JSON对象的文本形式） | select `'{"a":[1,2,3],"b":[4,5,6]}'::json#>>'{a,2}'`      | `3`            |

#### 例子1：

```
-- 简单标量/基本值
-- 基本值可以是数字、带引号的字符串、true、false或者null
SELECT '5'::json;

-- 简单标量/简单值,转化为jsonb类型
SELECT '5'::jsonb;

-- 有零个或者更多元素的数组（元素不需要为同一类型）
SELECT '[1, 2, "foo", null]'::json;

-- 包含键值对的对象
-- 注意对象键必须总是带引号的字符串
SELECT '{"bar": "baz", "balance": 7.77, "active": false}'::json;

-- 数组和对象可以被任意嵌套
SELECT '{"foo": [true, "bar"], "tags": {"a": 1, "b": null}}'::json;

-- "->" 通过键获得 JSON 对象域 结果为json对象
select '{"nickname": "goodspeed", "avatar": "avatar_url", "tags": ["python", "golang", "db"]}'::json->'nickname' as nickname;
 nickname
-------------
 "goodspeed"

-- "->>" 通过键获得 JSON 对象域 结果为text 
select '{"nickname": "goodspeed", "avatar": "avatar_url", "tags": ["python", "golang", "db"]}'::json->>'nickname' as nickname;
 nickname
-----------
 goodspeed
 
-- "->" 通过键获得 JSON 对象域 结果为json对象
select '{"nickname": "goodspeed", "avatar": "avatar_url", "tags": ["python", "golang", "db"]}'::jsonb->'nickname' as nickname;
 nickname
-------------
 "goodspeed"

-- "->>" 通过键获得 JSON 对象域 结果为text 
select '{"nickname": "goodspeed", "avatar": "avatar_url", "tags": ["python", "golang", "db"]}'::jsonb->>'nickname' as nickname;
 nickname
-----------
 goodspeed
```

#### 例子2：

```
CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  client TEXT NOT NULL,
  data JSONb NOT NULL
);

INSERT INTO books(client, data) values ( 'Joe',
    '{ "title": "Siddhartha", "author": { "first_name": "Herman", "last_name": "Hesse" } }'
),( 'Jenny',
    '{ "title": "Dharma Bums", "author": { "first_name": "Jack", "last_name": "Kerouac" } }'
),( 'Jenny',
    '{ "title": "100 años de soledad", "author": { "first_name": "Gabo", "last_name": "Marquéz" } }'
);

SELECT client,
   data->'title' AS title, data->'author' || '["c", "d"]'::jsonb AS author
   FROM books;

SELECT client,
   data->>'title' AS title, data->'author'->>'first_name' AS author
   FROM books;
```

### jsonb 额外的操作符

| 操作符 | 右操作数类型 | 描述                                                         | 例子                                                |
| ------ | :----------: | ------------------------------------------------------------ | --------------------------------------------------- |
| `@>`   |   `jsonb`    | Does the left JSON value contain the right JSON path/value entries at the top level?（左侧 JSON 值是否在顶层包含正确的 JSON 路径/值条目） | `'{"a":1, "b":2}'::jsonb @> '{"b":2}'::jsonb`       |
| `<@`   |   `jsonb`    | Are the left JSON path/value entries contained at the top level within the right JSON value?（左侧 JSON 路径/值条目是否包含在右侧 JSON 值的顶层） | `'{"b":2}'::jsonb <@ '{"a":1, "b":2}'::jsonb`       |
| `?`    |    `text`    | Does the *string* exist as a top-level key within the JSON value?（字符串是否作为 JSON 值中的顶级键存在） | `'{"a":1, "b":2}'::jsonb ? 'b'`                     |
| `?|`   |   `text[]`   | Do any of these array *strings* exist as top-level keys?（这些数组字符串中的任何一个是否作为顶级键存在） | `'{"a":1, "b":2, "c":3}'::jsonb ?| array['b', 'c']` |
| `?&`   |   `text[]`   | Do all of these array *strings* exist as top-level keys?（所有这些数组字符串是否都作为顶级键存在） | `'["a", "b"]'::jsonb ?& array['a', 'b']`            |
| `||`   |   `jsonb`    | Concatenate two `jsonb` values into a new `jsonb` value（将两个 jsonb 值连接成一个新的 jsonb 值） | `'["a", "b"]'::jsonb || '["c", "d"]'::jsonb`        |
| `-`    |    `text`    | Delete key/value pair or *string* element from left operand. Key/value pairs are matched based on their key value.（从左操作数中删除键/值对或字符串元素。键/值对根据其键值进行匹配） | `'{"a": "b"}'::jsonb - 'a'`                         |
| `-`    |   `text[]`   | Delete multiple key/value pairs or *string* elements from left operand. Key/value pairs are matched based on their key value.（从左操作数中删除多个键/值对或字符串元素。键/值对根据其键值进行匹配） | `'{"a": "b", "c": "d"}'::jsonb - '{a,c}'::text[]`   |
| `-`    |  `integer`   | Delete the array element with specified index (Negative integers count from the end). Throws an error if top level container is not an array.（删除具有指定索引的数组元素（负整数从末尾开始计数）。如果顶级容器不是数组，则会引发错误） | `'["a", "b"]'::jsonb - 1`                           |
| `#-`   |   `text[]`   | Delete the field or element with specified path (for JSON arrays, negative integers count from the end)（删除指定路径的字段或元素（对于JSON数组，负整数从末尾开始计数）） | `'["a", {"b":1}]'::jsonb #- '{1,b}'`                |
| `@?`   |  `jsonpath`  | Does JSON path return any item for the specified JSON value?（JSON 路径是否返回指定 JSON 值的任何项目） | `'{"a":[1,2,3,4,5]}'::jsonb @? '$.a[*] ? (@ > 2)'`  |
| `@@`   |  `jsonpath`  | Returns the result of JSON path predicate check for the specified JSON value. Only the first item of the result is taken into account. If the result is not Boolean, then `null` is returned.（返回指定 JSON 值的 JSON 路径谓词检查结果。只考虑结果的第一项。如果结果不是布尔值，则返回 null。） | `'{"a":[1,2,3,4,5]}'::jsonb @@ '$.a[*] > 2'`        |

例子：

```
-- @> 左侧 JSON 值是否在顶层包含正确的 JSON 路径/值条目
select '{"a":1, "b":2}'::jsonb @> '{"b":2}'::jsonb

-- <@ 左侧 JSON 路径/值条目是否包含在右侧 JSON 值的顶层
select '{"b":2}'::jsonb <@ '{"a":1, "b":2}'::jsonb

-- ? 字符串是否作为 JSON 值中的顶级键存在
select '{"a":1, "b":2}'::jsonb ? 'b'

-- ?| 这些数组字符串中的任何一个是否作为顶级键存在
select '{"a":1, "b":2, "c":3}'::jsonb ?| array['b', 'c']

-- ?& 所有这些数组字符串是否都作为顶级键存在
select '["a", "b"]'::jsonb ?& array['a', 'b']
-- true

-- || 将两个 jsonb 值连接成一个新的 jsonb 值
select '["a", "b"]'::jsonb || '["c", "d"]'::jsonb
--["a", "b", "c", "d"]

-- - 从左操作数中删除键/值对或字符串元素。键/值对根据其键值进行匹配
select '{"a": "b"}'::jsonb - 'a'

-- - 从左操作数中删除多个键/值对或字符串元素。键/值对根据其键值进行匹配
select '{"a": "b", "c": "d"}'::jsonb - '{a,c}'::text[]
-- {}

-- - 删除具有指定索引的数组元素（负整数从末尾开始计数）。如果顶级容器不是数组，则会引发错误
select '["a", "b"]'::jsonb - 1
-- ["a"]

-- #- 删除指定路径的字段或元素（对于JSON数组，负整数从末尾开始计数）
select '["a", {"b":1}]'::jsonb #- '{1,b}'
--["a", {}]

-- @? JSON 路径是否返回指定 JSON 值的任何项目
select '{"a":[1,2,3,4,5]}'::jsonb @? '$.a[*] ? (@ > 2)'
-- true

-- @@ 返回指定 JSON 值的 JSON 路径谓词检查结果。只考虑结果的第一项。如果结果不是布尔值，则返回 null。
select '{"a":[1,2,3,4,5]}'::jsonb @@ '$.a[*] > 2'
-- true
```

### JSON 方法

#### 删除某个字段

```
update books set data = data - 'nickname' where client='Joe';
```

#### jsonb_set

```
jsonb_set(target jsonb, path text[], new_value jsonb [, create_missing boolean])
参数：
    target：目标（jsonb类型的属性）
    path ：路径，如果jsonb是数组‘{0，a}’表示在下标是0的位置更新a属性，如果不是数组，是对象，则写‘{a}’即可
    new_value：新值
选填参数：
	create_missing：jsonb字段不存在时创建，默认为true
返回：
	更新后的jsonb
例子：
	1.jsonb_set('{"title": {"Id":"abcd","fullName": "小A"}}','{title,fullName}','"小B"');
	2.update books set data = jsonb_set(data, '{nickname}', '"gs1"', false) where client='Joe';
```

#### jsonb_insert

```
jsonb_insert(target jsonb, path text[], new_value jsonb [, insert_after boolean])
	返回插入 new_value 的目标。如果 path 指定的 target 部分在 JSONB 数组中，如果 insert_after 为 true（默认为 false），则 new_value 将插入 target 之前或之后。如果 path 指定的 target 部分在 JSONB 对象中，则只有在 target 不存在时才会插入 new_value。与面向路径的运算符一样，出现在路径中的负整数从 JSON 数组的末尾开始计数。
例子：
	1. update books set data = jsonb_insert(data, '{nickname,0}', '{"a": "a1"}', false) where client='Joe';
```

#### jsonb_pretty

```
jsonb_pretty(from_json jsonb)
返回json格式
例子：
SELECT client,
   data->>'title' AS title, jsonb_pretty(data) AS author
   FROM books;
Ouput:
_________________________________________________________
client|   title     | author                             |
______|_____________|____________________________________|
Jenny | Dharma Bums | "{
      |             |       ""title"": ""Dharma Bums"",
      |             |      ""author"": {
      |             |          ""last_name"": ""Kerouac"",
      |             |           ""first_name"": ""Jack""
      |             |       }
      |             |   }"Jenny,Dharma | Bums,"{
      |             |       ""title"": ""Dharma Bums"",
      |             |      ""author"": {
      |             |          ""last_name"": ""Kerouac"",
      |             |           ""first_name"": ""Jack""
      |             |       }
      |             |   }"

```

还有许多json方法请看文档。

### jsonb 索引

- gin索引
- btree索引
- `hash`索引

​	PostgreSQL中通过GIN索引则实现了索引的**模式自由**，即索引时不需要指定JSON数据中的指定字段，后续就可以按jsonb中已的任意的字段进行查询，即对jsonb里面的字段都建立索引。而使用btree索引，这是就不是**模式自由**，btree索引是使用sonb里面的某一个字段建立索引。




### Reference

- https://www.postgresql.org/docs/current/datatype-json.html
- https://www.postgresql.org/docs/12/functions-json.html#FUNCTIONS-JSON-PROCESSING
- http://www.pgsql.tech/article_104_10000050
- http://www.postgres.cn/docs/10/datatype-json.html
