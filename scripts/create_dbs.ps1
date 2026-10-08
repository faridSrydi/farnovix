$createdb = "C:\Program Files\PostgreSQL\18\bin\createdb.exe"
$psql = "C:\Program Files\PostgreSQL\18\bin\psql.exe"

$dbs = @("user_db", "event_db", "booking_db", "payment_db", "ticket_db", "notification_db")
foreach ($db in $dbs) {
    Write-Output "Creating database: $db ..."
    & $createdb -h 127.0.0.1 -p 5433 -U ticketing $db 2>&1
}

Write-Output "Listing all databases:"
& $psql -h 127.0.0.1 -p 5433 -U ticketing -d postgres -c "\l"
