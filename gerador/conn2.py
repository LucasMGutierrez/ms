import sys
import operator
from cassandra.cluster import Cluster

class Span:

    def __init__(self, span_id, start_time, duration, trace_id, operation):
        self.span_id = span_id
        self.start_time = start_time
        self.duration = duration
        self.end_time = start_time + duration
        self.trace_id = trace_id
        self.operation = operation

class Service:

    def __init__(self, name, trace_ids):
        self.name = name
        self.childof = None
        self.spans = []
        self.childs = []
        self.demand = 0.0
	self.v = 0
        self.trace_ids = trace_ids
        self.root = False
        self.operation = None
        self.parallel = set()

    def validoperation(self, span):
        return self.operation == None or self.operation == span.operation

    def setparent(self, service):
        self.childof = service

    def setparent2(self, service):
        if self != service:
	    self.childof = service
	    self.v += 1

    def setroot(self, b):
        self.root = b

    def isroot(self):
        return self.root

    def addspan(self, span):
        self.spans.append(span)

    def numspans(self):
        return len(self.spans)

    def addchild(self, span):
        self.childs.append(span)

    def haschild(self):
        return self.childs != []

    def meanspans(self):
        mean = 0.0
        cont = 0.0

        for span in self.spans:
            if self.validoperation(span):
                cont += 1.0
                mean += span.duration

        return mean / cont

    # TODO numero de visitas
    def setdemand(self):
        self.demand = self.meanspans()

        for child in self.childs:
            child.setdemand()
            self.demand -= child.meanspans() * child.visits()

    def visits(self):
        if self.isroot():
            return 1
        return self.v / len(self.trace_ids)

    def adjust(self):
        if self.visits() == 1 and len(self.trace_ids) != len(self.spans):
            bigger = 0
            for trace_id in self.trace_ids:
                for span in self.spans:
                    if span.trace_id == trace_id and span.duration > bigger:
                        self.operation = span.operation
                        bigger = span.duration

    def sameinterval(self, s1, s2):
        return (s1.end_time < s2.end_time and s1.end_time > s2.start_time) or (s2.end_time < s1.end_time and s2.end_time > s1.start_time)

    def setparallel(self):
        for c1 in self.childs:
            for c2 in self.childs:
                for s1 in c1.spans:
                    if c1.validoperation(s1):
                        for s2 in c2.spans:
                            if c2.validoperation(s2) and s1.span_id != s2.span_id and self.sameinterval(s1, s2):
                                self.parallel.add(c1)
                                self.parallel.add(c2)

# TODO class trace ??
class Trace:
    pass

# Get a service operation name
def getOperation(session, any_service):
    query = 'select operation_name from operation_names where service_name=\'' + any_service + '\';';
    rows = session.execute(query)

    return rows[0].operation_name

# Get all trace_id
def getAllTraceId(session, any_service, operation):
    query = 'select trace_id from service_operation_index where service_name = \'' + any_service + '\' and operation_name = \'' + operation + '\';'

    rows = session.execute(query)

    trace_ids = []
    for row in rows:
        trace_ids.append('0x{}'.format(row.trace_id.encode('hex')))

    return trace_ids

# Get all service names
# TODO considerar todos os traces
def getServiceNames(session, trace_ids):
    query = 'select service_name from service_name_index where trace_id = ' +  trace_ids[0] + ' ALLOW FILTERING;'

    rows = session.execute(query)

    service_names = []
    for row in rows:
        service_names.append(row.service_name)

    return service_names

# Get all service names
def getServiceNames2(session, trace_ids):
    for trace_id in trace_ids:
        query = 'select service_name from service_name_index where trace_id = ' +  trace_id + ' ALLOW FILTERING;'
        rows = session.execute(query)
        
        service_names = set()
        for row in rows:
            service_names.add(row.service_name)

    return service_names

# Create Services
def createServices(service_names, trace_ids):
    services = dict()

    for service_name in service_names:
        services[service_name] = Service(service_name, trace_ids)

    return services

# Get spans from traces
def getSpans(session, trace_ids, services):
    for trace_id in trace_ids:
        query = 'select * from traces where trace_id = ' + trace_id + ';'
        rows = session.execute(query)

        for row in rows:
            span = Span(row.span_id, row.start_time, row.duration, trace_id, row.operation_name)
            services[row.process.service_name].addspan(span)

# Define childOf
# TODO todos os filhos, em todos os traces
def setChildOf(session, services, trace_ids):
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

# Define childOf
def setChildOf2(session, services, trace_ids):
    for trace_id in trace_ids:
        span_service = dict()
        query = 'select * from traces where trace_id = ' + trace_id + ';'
        rows = session.execute(query)
        for row in rows:
            span_service[row.span_id] = row.process.service_name
    
        query = 'select * from traces where trace_id = ' + trace_id + ';'
        rows = session.execute(query)
        for row in rows:
            # print row.refs
            if row.refs != None and row.refs[0].span_id != 0:
                services[span_service[row.span_id]].setparent2(services[span_service[row.refs[0].span_id]])

# Define childs
def setChilds(session, services, services_names):
    for service_name in service_names:
        if services[service_name].childof != None:
            services[service_name].childof.addchild(services[service_name])

# Find Root Service
def findRootService(services, services_names):
    for service_name in service_names:
        if services[service_name].isroot():
            root = services[service_name]

    return root

def adjustServices(services, service_names):
    for s in service_names:
        services[s].adjust()

def checkParallelism(services, service_names):
    for s in service_names:
        services[s].setparallel()

########################################
#                                      #
#                SCRIPT                #
#                                      #
########################################

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

operation = getOperation(session, any_service)
trace_ids = getAllTraceId(session, any_service, operation)
service_names = getServiceNames2(session, trace_ids)
services = createServices(service_names, trace_ids)
getSpans(session, trace_ids, services)
setChildOf2(session, services, trace_ids)
setChilds(session, services, service_names)
#root = findRootService(services, services_names)
root = services[any_service]
root.setroot(True)
adjustServices(services, service_names)
checkParallelism(services, service_names)

for s in service_names:
    print s, services[s].operation, services[s].visits()
    for i in services[s].parallel:
        print '\t' + i.name

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
f.write('SetComment("Modelo de Redes de Filas para o ' + app + '")\n\n')

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
    demand = str(services[service_name].demand / 1000.0)
    visits = str(services[service_name].visits())
    f.write('SetDemand("' + service_name + '", "req", ' + visits + ' * ' + demand + ')\n')
f.write('\n')

f.write('# Soluciona o Modelo\n\n')

f.write('Solve(CANON)\n\n')

f.write('Report()\n\n')
