import os

def removeall(l, item):
	while item in l:
		l.remove(item)

def microservice(str):
	return {
		'ms1': True,
		'ms2': True,
		'ms3': True,
		'frontend': True,
	}.get(str, False)

with open('arq') as f:
	lines = f.readlines()

for l in lines:
	ret = l.rsplit(' ')
	removeall(ret, '')
	if microservice(ret[0]):
		os.system("kill " + ret[1])