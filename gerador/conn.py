import sys
from cassandra.cluster import Cluster

class Span:

    def __init__(self, span_id, start_time, duration):
        self.span_id = span_id
        self.start_time = start_time
        self.duration = duration

class Service:

    def __init__(self, name):
        self.name = name
        self.childof = None
        self.spans = []
        self.childs = []
        self.demand = 0.0

    def setparent(self, service):
        self.childof = service

    def isroot(self):
        return self.childof == None

    def addspan(self, span):
        self.spans.append(span)

    def numspans(self):
        return len(self.spans)

    def addchild(self, span):
        self.childs.append(span)

    def meanspans(self):
        mean = 0.0
        for span in self.spans:
            mean += span.duration

        return mean / len(self.spans)

    def setdemand(self):
        self.demand = self.meanspans()

        for child in self.childs:
            child.setdemand()
            self.demand -= child.meanspans()

# Check args
if len(sys.argv) == 1:
    print 'No arguments'
    exit()

any_service = sys.argv[1]

if len(sys.argv) == 2:
    app = any_service
else:
    app = sys.argv[2]

# Cassandra
cluster = Cluster()
session = cluster.connect('jaeger_v1_dc1')

# Get a service operation name
query = 'select operation_name from operation_names where service_name=\'' + any_service + '\';';

rows = session.execute(query)
operation = rows[0].operation_name

# Get all trace_id
query = 'select trace_id from service_operation_index where service_name = \'' + any_service + '\' and operation_name = \'' + operation + '\';'

rows = session.execute(query)

trace_ids = []
for row in rows:
    trace_ids.append('0x{}'.format(row.trace_id.encode('hex')))

# Get all service names
query = 'select service_name from service_name_index where trace_id = ' +  trace_ids[0] + ' ALLOW FILTERING;'

rows = session.execute(query)

service_names = []
for row in rows:
    service_names.append(row.service_name)

# Create Services
services = dict()

for service_name in service_names:
    services[service_name] = Service(service_name)

# Get spans from traces
for trace_id in trace_ids:
    query = 'select * from traces where trace_id = ' + trace_id + ';'
    rows = session.execute(query)

    for row in rows:
        span = Span(row.span_id, row.start_time, row.duration)
        services[row.process.service_name].addspan(span)

# Define childOf
span_service = dict()
query = 'select * from traces where trace_id = ' + trace_ids[0] + ';'
rows = session.execute(query)
for row in rows:
    span_service[row.span_id] = row.process.service_name

query = 'select * from traces where trace_id = ' + trace_ids[0] + ';'
rows = session.execute(query)
for row in rows:
    if row.refs != None:
        services[span_service[row.span_id]].setparent(services[span_service[row.refs[0].span_id]])

# Define childs
for service_name in service_names:
    if not services[service_name].isroot():
        services[service_name].childof.addchild(services[service_name])

# Find Root Service
for service_name in service_names:
    if services[service_name].isroot():
        root = services[service_name]

# Calculate Demands
root.setdemand()

# Create R file
f = open(app + '.R', 'w')

f.write('\n# Modelo de Redes de Filas para o ' + app + '\n\n')
f.write('library(pdq)\n\n')
f.write('nReqs <- 1\n')
f.write('tempo <- ' + str(2.0 * root.meanspans() / 1000.0) + '\n')
f.write('lambda <- nReqs / tempo\n\n')

f.write('# Inicializacao\n\n')

f.write('Init("OpenCircuit")\n')
f.write('SetComment("Mode de Redes de Filas para o ' + app + '")\n\n')

f.write('# Definicao das Filas\n\n')

for service_name in service_names:
    f.write('CreateNode("' + service_name + '", CEN, FCFS)\n')
f.write('\n')

f.write('# Definicao da Carga de Trabalho\n\n')

f.write('CreateOpen("req", lambda)\n\n')
f.write('SetWUnit("req")\n')
f.write('SetTUnit("ms")\n\n')

f.write('# Definicao das Demandas\n\n')

for service_name in service_names:
    f.write('SetDemand("' + service_name + '", "req", ' + str(services[service_name].demand / 1000.0) + ')\n')
f.write('\n')

f.write('# Soluciona o Modelo\n\n')

f.write('Solve(CANON)\n\n')

f.write('Report()\n\n')

