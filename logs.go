package main

import (
	"time"

	"flag"
	"fmt"
	"os"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var (
	sess *session.Session
	svc  *cloudwatchlogs.CloudWatchLogs
)

func main() {
	start := flag.Int64("start", 0, "start unix time (seconds)")
	end := flag.Int64("end", 0, "end unix time (seconds)")
	filter := flag.String("filter", "", "Filter pattern (checkout http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)")
	outputFile := flag.String("f", "", "Output file name (optional)")
	region := flag.String("region", "eu-west-1", "AWS Region")
	logGroupName := flag.String("log-group-name", "", "Log Group Name")
	flag.Parse()

	sess = session.New(aws.NewConfig().WithRegion(*region))
	svc = cloudwatchlogs.New(sess)

	if *start == 0 && *end == 0 {
		*end = time.Now().Unix()
		*start = *end - 120
	}

	req := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:     aws.Int64(*start * 1e3),
		Interleaved:   aws.Bool(true),
		LogGroupName:  aws.String(*logGroupName),
		EndTime:       aws.Int64(*end * 1e3),
		FilterPattern: aws.String(*filter),
	}
	fmt.Println("Timespan", time.Unix(*start, 0), time.Unix(*end, 0))
	if *outputFile == "" {
		*outputFile = fmt.Sprintf("%s_%s_%s.log", *logGroupName, time.Unix(*start, 0).Format(time.RFC3339), time.Unix(*end, 0).Format(time.RFC3339))
	}
	file, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed while creating output file")
	}
	defer file.Close()

	var earliest, latest int64
	earliest = 1 << 62
	for {
		res, err := svc.FilterLogEvents(req)
		if ae, ok := err.(awserr.Error); ok && ae.Code() == "ThrottlingException" {
			// backoff
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Println(err)
			continue
		}

		for _, ev := range res.Events {
			if ev.Message != nil && ev.Timestamp != nil {

				stdTs := int64(float64(*ev.Timestamp) / 1000)
				if stdTs < earliest {
					earliest = stdTs
				}

				if stdTs > latest {
					latest = stdTs
				}
				_, err := file.WriteString("[" + time.Unix(stdTs, 0).String() + "]" + " " + *ev.Message + "\n")
				if err != nil {
					log.Println(err)
				}

				stats, _ := file.Stat()
				str := fmt.Sprintf("Output size: %.2f M, processed logs from: %v, to: %v", float64(stats.Size())/float64((1<<20)), time.Unix(earliest, 0), time.Unix(latest, 0))
				os.Stdout.Write([]byte("\r" + str))
				os.Stdout.Sync()
			}
		}

		if res.NextToken == nil {
			return
		}
		req.NextToken = res.NextToken
	}
}
