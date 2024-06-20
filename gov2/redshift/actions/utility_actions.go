package actions

//// snippet-start:[gov2.redshift.Movie.struct]
//
//// Movie makes it easier to use Movie objects given in json format.
//type Movie struct {
//	ID    int    `json:"id"`
//	Title string `json:"title"`
//	Year  int    `json:"year"`
//}
//
//// snippet-end:[gov2.redshift.Movie.struct]
//
//// snippet-start:[gov2.redshift.RedshiftQuery.struct]
//
//// RedshiftQuery makes it easier to deal with RedshiftQuery objects.
//type RedshiftQuery struct {
//	Result interface{}
//	Input  redshiftdata.DescribeStatementInput
//	Client *redshiftdata.Client
//}
//
//// snippet-end:[gov2.redshift.RedshiftQuery.struct]
//
//// snippet-start:[gov2.redshift.WaitForQueryStatus]
//
//// WaitForQueryStatus waits until the given RedshiftQuery object has succeeded or failed.
//func WaitForQueryStatus(query RedshiftQuery, showProgress bool) error {
//	done := false
//	attempts := 0
//	maxWaitCycles := 30
//	for done == false {
//		describeOutput, err := query.Client.DescribeStatement(context.TODO(), &query.Input)
//		if err != nil {
//			return err
//		}
//		if describeOutput.Status == "FAILED" {
//			return errors.New("failed to describe statement")
//		}
//		if attempts >= maxWaitCycles {
//			return errors.New("timed out waiting for statement")
//		}
//		if showProgress {
//			fmt.Print(".")
//		}
//		if describeOutput.Status == "FINISHED" {
//			done = true
//		}
//		attempts++
//	}
//	return nil
//}
//
//// snippet-end:[gov2.redshift.WaitForQueryStatus]
//
//// snippet-start:[gov2.redshift.loadMoviesFromJSON]
//
//// loadMoviesFromJSON takes the "Movies.json" file and populates a slice of Movie objects.
//func loadMoviesFromJSON(filename string, filesystem demotools.IFileSystem) ([]Movie, error) {
//	file, err := filesystem.Open("../../resources/sample_files/" + filename)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var movies []Movie
//	err = json.NewDecoder(file).Decode(&movies)
//	if err != nil {
//		return nil, err
//	}
//
//	return movies, nil
//}
//
//// snippet-end:[gov2.redshift.loadMoviesFromJSON]
//
//// snippet-start:[gov2.redshift.WaitForClusterAvailable]
//
//// WaitForClusterAvailable loops until a success or failure state is returned about the given cluster.
//func WaitForClusterAvailable(client *redshift.Client, clusterId string, questioner demotools.IQuestioner) {
//	questioner.Ask("Now we will wait until the cluster is available. Press Enter to continue...")
//
//	log.Println("Waiting for cluster to become available. This may take a few minutes.")
//	describeClusterInput := &redshift.DescribeClustersInput{
//		ClusterIdentifier: &clusterId,
//	}
//
//	for {
//		output, err := client.DescribeClusters(context.TODO(), describeClusterInput)
//		if err != nil {
//			log.Printf("Failed to describe cluster: %v\n", err)
//			return
//		}
//
//		if output.Clusters[0].ClusterStatus != nil && *output.Clusters[0].ClusterStatus == "available" {
//			log.Println("Cluster is available! Total Elapsed Time:", time.Since(time.Now().Add(-1*time.Second*time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
//			break
//		}
//
//		log.Print("Elapsed Time: ")
//		log.Println(time.Since(time.Now().Add(-1 * time.Second * time.Duration(output.Clusters[0].ClusterCreateTime.Second()))))
//		// TODO implement demotools Pause
//		time.Sleep(5 * time.Second)
//	}
//}
//
//// snippet-end:[gov2.redshift.WaitForClusterAvailable]
//
//// snippet-start:[gov2.redshift.buildSqlStatements]
//
//// buildSqlStatements formats a series of SQL statements to be batch fed to the Redshift cluster.
//func buildSqlStatements(movies []Movie, numRecords int) []string {
//	var sqlStatements []string
//
//	for i, movie := range movies {
//		if i >= numRecords {
//			break
//		}
//
//		sqlStatement := fmt.Sprintf("INSERT INTO Movies (title, year) VALUES ('%s', %d);", movie.Title, movie.Year)
//		sqlStatements = append(sqlStatements, sqlStatement)
//	}
//
//	return sqlStatements
//}
//
//// snippet-end:[gov2.redshift.buildSqlStatements]
//
//// snippet-start:[gov2.redshift.PopulateMoviesTable]
//
//// PopulateMoviesTable reads data from the "Movies.json" file and inserts records into the "Movies" table.
//func PopulateMoviesTable(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
//	fmt.Println("Populate the Movies table using the Movies.json file.")
//	numRecords := questioner.AskInt(
//		fmt.Sprintf("Enter a value between %v and %v:", 10, 100),
//		demotools.InIntRange{Lower: 10, Upper: 100})
//
//	movies, err := loadMoviesFromJSON("Movies.json")
//	if err != nil {
//		fmt.Printf("Failed to load movies from JSON: %v\n", err)
//		return
//	}
//
//	sqlStatements := buildSqlStatements(movies, numRecords)
//
//	input := &redshiftdata.BatchExecuteStatementInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//		Sqls:              sqlStatements,
//	}
//
//	result, err := client.BatchExecuteStatement(context.TODO(), input)
//	if err != nil {
//		fmt.Printf("Failed to execute batch statement: %v\n", err)
//		return
//	}
//
//	describeInput := redshiftdata.DescribeStatementInput{
//		Id: result.Id,
//	}
//
//	query := RedshiftQuery{
//		Client: client,
//		Result: result,
//		Input:  describeInput,
//	}
//	err = WaitForQueryStatus(query, true)
//	if err != nil {
//		fmt.Printf("Failed to execute batch insert query: %v\n", err)
//		return
//	}
//	fmt.Printf("Successfully executed batch statement\n")
//
//	fmt.Printf("%d records were added to the Movies table.\n", numRecords)
//}
//
//// snippet-end:[gov2.redshift.PopulateMoviesTable]
//
//// snippet-start:[gov2.redshift.DeleteDataRows]
//
//// DeleteDataRows deletes all rows from the "Movies" table.
//func DeleteDataRows(client *redshiftdata.Client, userName, clusterId string) (bool, error) {
//	deleteRows := &redshiftdata.ExecuteStatementInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//		Sql:               aws.String("DELETE FROM Movies;"),
//	}
//
//	result, err := client.ExecuteStatement(context.TODO(), deleteRows)
//	if err != nil {
//		fmt.Printf("Failed to execute batch statement: %v\n", err)
//		return false, err
//	}
//	describeInput := redshiftdata.DescribeStatementInput{
//		Id: result.Id,
//	}
//	query := RedshiftQuery{
//		Client: client,
//		Result: result,
//		Input:  describeInput,
//	}
//	err = WaitForQueryStatus(query, true)
//	if err != nil {
//		fmt.Printf("Failed to execute delete query: %v\n", err)
//		return false, err
//	}
//
//	fmt.Printf("Successfully executed delete statement\n")
//	return true, nil
//}
//
//// snippet-end:[gov2.redshift.DeleteDataRows]
//
//// snippet-start:[gov2.redshift.QueryMoviesByYear]
//
//// QueryMoviesByYear retrieves only movies from the "Movies" table which match the given year.
//func QueryMoviesByYear(client *redshiftdata.Client, userName, clusterId string, questioner demotools.IQuestioner) {
//	fmt.Println("Query the Movies table by year.")
//	year := questioner.AskInt(
//		fmt.Sprintf("Enter a value between %v and %v:", 2012, 2014),
//		demotools.InIntRange{Lower: 2012, Upper: 2014})
//
//	input := &redshiftdata.ExecuteStatementInput{
//		ClusterIdentifier: &clusterId,
//		Database:          aws.String("dev"),
//		DbUser:            &userName,
//		Sql:               aws.String(fmt.Sprintf("SELECT title FROM Movies WHERE year = %d;", year)),
//	}
//
//	result, err := client.ExecuteStatement(context.TODO(), input)
//	if err != nil {
//		fmt.Printf("Failed to query movies: %v\n", err)
//		return
//	}
//
//	fmt.Println("The identifier of the statement is", *result.Id)
//
//	describeInput := redshiftdata.DescribeStatementInput{
//		Id: result.Id,
//	}
//
//	query := RedshiftQuery{
//		Client: client,
//		Input:  describeInput,
//		Result: result,
//	}
//	err = WaitForQueryStatus(query, true)
//	if err != nil {
//		fmt.Printf("Failed to execute query: %v\n", err)
//	}
//	fmt.Printf("Successfully executed query\n")
//
//	getResultOutput, err := client.GetStatementResult(context.TODO(), &redshiftdata.GetStatementResultInput{
//		Id: result.Id,
//	})
//	if err != nil {
//		fmt.Printf("Failed to query movies: %v\n", err)
//		return
//	}
//	for _, row := range getResultOutput.Records {
//		for _, col := range row {
//			title, ok := col.(*types.FieldMemberStringValue)
//			if !ok {
//				fmt.Println("Failed to parse the field")
//			} else {
//				fmt.Printf("The Movie title field is %s\n", title.Value)
//			}
//		}
//	}
//}
//
//// snippet-end:[gov2.redshift.QueryMoviesByYear]
