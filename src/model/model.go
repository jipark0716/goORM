package model

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
    "strings"
    "config"
)

type Query struct {
    where string
    column []string
    limit int
    offset int
    table string
    order []string
    connection string
}

type Example struct {
    Idx int `db:"id"`
    Msg string `db:"msg"`
}

var DB map[string]*sqlx.DB

/*
 * select query 실행
 *
 * @param
 *     interface $dest 결과받을 변수
 *
 * @return void
*/
func (query *Query) Get(dest interface{}) {
    e := fmt.Sprintf("%#v", dest)
    e = strings.Replace(e, "&[]", "", 1)
    e = strings.Replace(e, "{}", "", 1)
    queryString := query.Sql(e)

    db := Conn(query.connection)
    err := db.Select(dest, queryString)
    if err != nil{
        println(queryString)
        fmt.Printf("%#v\n", err)
    }
}

/*
 * DB 연결
 *
 * @param
 *     string connection 연결DB 이름
 *
 * @return *sqlx.DB 연결된 DB
*/
func Conn(connection string) *sqlx.DB{
    if DB == nil{
        DB = make(map[string]*sqlx.DB)
    }
    db := DB[connection]
    if db != nil {
        return db
    }
    db, err := sqlx.Connect(
        "mysql",
        fmt.Sprintf(
            "%s:%s@(%s:3306)/%s",
            config.Get(fmt.Sprintf("DB_%s_USER", strings.ToUpper(connection))),
            config.Get(fmt.Sprintf("DB_%s_PASS", strings.ToUpper(connection))),
            config.Get(fmt.Sprintf("DB_%s_HOST", strings.ToUpper(connection))),
            config.Get(fmt.Sprintf("DB_%s_NAME", strings.ToUpper(connection))),
        ),
    )
    if err != nil {
        fmt.Println(err)
    }
    row, _ := db.Queryx("SET NAMES utf8mb4")
    defer row.Close()

    DB[connection] = db
    return db
}

/*
 * Limit 추가
 *
 * @param
 *     int $limit
 *
 * @return void
*/
func (query *Query) Limit(limit int) {
    query.limit = limit
}

/*
 * 정렬 추가
 *
 * @param
 *     string $col 정렬할 column
 *     string $sc 정렬 순서
 *
 * @return void
*/
func (query *Query) AddOrder(col string, sc string){
    query.order = append(query.order, fmt.Sprintf("`%s` %s", col, sc))
}

/*
 * where절 추가
 *
 * @param
 *     string $where Raw 조건문
 *
 * @return void
*/
func (query *Query) Where(where string){
    if query.where != "" {
        query.where += " and "
    }else{
        query.where += "where "
    }
    query.where += where
}

/*
 * 쿼리 문자열로 변환
 *
 * @param
 *     string $e model.struct name
 *
 * @return string 쿼리 문자열
*/
func (query *Query) Sql(param ...string) string{
    var e string
    if(len(param) > 0){
        e = param[0]
    }
    switch e {
    case "model.Example":
        query.table = "test"
        if len(query.column) == 0 {
            query.column = []string{"id", "msg"}
        }
    }
    col := strings.Join(query.column, ",")
    order := ""
    if len(query.order) > 0 {
        order =  "ORDER BY" + strings.Join(query.order, " ")
    }
    result := fmt.Sprintf("SELECT %s FROM `%s` %s %s", col, query.table, query.where, order)
    if query.limit > 0 {
        result += fmt.Sprintf(" LIMIT %d", query.limit)
    }
    return result
}

/*
 * 테이블 변경
 *
 * @param
 *     string $table 변경할 테이블 이름
 *
 * @return void
*/
func (query *Query) Table(table string) {
    query.table = table
}

/*
 * count(*)
 *
 * @return int
*/
func (query *Query) Count() int {
    query.column = []string{"count(*) as cnt"}
    queryString := query.Sql()
    db := Conn(query.connection)
    row := db.QueryRowx(queryString)
    if row == nil {
        println(queryString)
        return 0
    }
    var count int
    row.Scan(&count)

    return count
}

/*
 * 쿼리빌더 생성자
 *
 * @return model.Query
*/
func NewQuery() Query {
    return Query{where : "",connection:"master",limit:1}
}
