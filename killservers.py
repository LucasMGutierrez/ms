import os
import os.path

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

if not os.path.isfile('arq'):
    print('File does not exist. Execute:')
    print('sudo lsof -i -P -n | grep LISTEN > arq')
    exit()

with open('arq') as f:
    lines = f.readlines()

for l in lines:
    ret = l.rsplit(' ')
    removeall(ret, '')
    if microservice(ret[0]):
        print(ret[0])
        os.system("kill " + ret[1])

os.system('rm arq')
