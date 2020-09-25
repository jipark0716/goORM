package helper

import (
    "os"
    "fmt"
    "time"
)

/*
 * 로그남기기
 *
 * @param
 *     string $tag 태그
 *     string $content 내용
 *
 * @return void
*/
func Log(tag string, content string) {
    now := time.Now();
    os.MkdirAll("log", os.ModePerm)
    file, openErr := os.OpenFile(fmt.Sprintf("log/%s_%s.log", tag, now.Format("20060102")), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

    if openErr != nil {
        fmt.Printf("%#v\n", openErr);
        return
    }
    defer file.Close();

    _ , writeErr := file.Write([]byte(fmt.Sprintf("%s\n", content)))

    if writeErr != nil {
        fmt.Printf("%#v\n", writeErr);
    }
}

/*
 * 배열 분할하기
 *
 * @param
 *     []string $array 분할할 배열
 *     int $size
 *
 * @return [][]string
*/
func Chunk(array []string, size int) (chunks [][]string) {
    for size < len(array) {
        array, chunks = array[size:], append(chunks, array[0:size:size])
    }

    return append(chunks, array)
}

/*
 * 격주 짝수 홀수 확인
 *
 * @param
 *     time.Time 확인할 시간
 *
 * @return string biweeks:격주 짝수, biweek:격주 홀수
*/
func EOWeek(time time.Time) string {
    _, week := time.ISOWeek()
    if week % 2 == 0 {
        return "biweeks"
    } else {
        return "biweek"
    }
}
