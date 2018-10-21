
# Modelo de Redes de Filas para o MS

library(pdq)

nReqs <- 1
tempo <- 4784.1684
lambda <- nReqs / tempo

# Inicializacao

Init("OpenCircuit")
SetComment("Modelo de Redes de Filas para o MS")

# Definicao das Filas

CreateNode("Ms3", CEN, FCFS)
CreateNode("Ms2", CEN, FCFS)
CreateNode("Frontend", CEN, FCFS)
CreateNode("Hello-World", CEN, FCFS)
CreateNode("Ms1", CEN, FCFS)

# Definicao da Carga de Trabalho

CreateOpen("req", lambda)

SetWUnit("req")
SetTUnit("ms")

# Definicao das Demandas

SetDemand("Ms3", "req", 991.9154)
SetDemand("Ms2", "req", 998.7504)
SetDemand("Frontend", "req", 100.2104)
SetDemand("Hello-World", "req", 1.1358)
SetDemand("Ms1", "req", 300.0722)

# Soluciona o Modelo

Solve(CANON)

Report()

