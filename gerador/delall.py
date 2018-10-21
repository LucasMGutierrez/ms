from cassandra.cluster import Cluster

# Cassandra
cluster = Cluster()
session = cluster.connect('jaeger_v1_dc1')

query = 'TRUNCATE dependencies'
session.execute(query)

query = 'TRUNCATE duration_index'
session.execute(query)

query = 'TRUNCATE operation_names'
session.execute(query)

query = 'TRUNCATE service_name_index'
session.execute(query)

query = 'TRUNCATE service_names'
session.execute(query)

query = 'TRUNCATE service_operation_index'
session.execute(query)

query = 'TRUNCATE tag_index'
session.execute(query)

query = 'TRUNCATE traces'
session.execute(query)

