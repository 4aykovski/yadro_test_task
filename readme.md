# YADRO Test Task
��������� �������� �������� �������, ������� ������ �� ������� �������������
�����, ������������ ������� � ������������ ������� �� ���� � ����� ���������
������� �����.
������� ����� ���� ����������� �� Golang.

�������� ������� �����: ���� ��� ��������� ������ � �������� ����� ���������
�� ����� Golang (������ 1.19 � ������) � �������������� go modules, ���������� ��
������� � �������� ������� (���������� ������ � �� ���������� ������������).
������� ������ ������������ ����� ��������� ����. ���� ����������� ������
���������� ��� ������� ���������. ������ ������� ���������:

$ task.exe test_file.txt

��������� ������ ���������c� � Linux ��� Windows � �������������� docker
container-a (��������� ��������� Dockerfile). ��������� ������������� �����������
���������� (https://pkg.go.dev/std). ������������� ����� ��������� ���������, �����
�����������, ���������. � �������, ����� ������ � �������� �����, ���������
������������ ���������� �� ������� ��������� ��� ��������.

# ������ ��� ���������� ����� ����-������

1. ����������� �����������
```
    git clone https://github.com/4aykovski/yadro_test_task.git
```
2.  ������� � ���������� �������
```
    cd yadro_test_task
```
3. ��������� �����
```
    docker build -t your_image_name:your_tag .
```
4. ��������� ��������� � ��������������� ����-�������
```
    docker run -t your_image_name:your_tag
```

# ������ � ����������� ����� ����-������
1. ����������� �����������
```
    git clone https://github.com/4aykovski/yadro_test_task.git
```
2.  ������� � ���������� �������
```
    cd yadro_test_task
```
3. ������� � ����� `cases` ����� ����-���� � ��������� ���
```
    nano ./cases/new_test_case.txt
```
4. �� ������� �������� � Makefile ����� ����-���� �� �������� � ������������� ��� ������� ���� ����-������ �����
```
    nano ./Makefile
```
5. ��������� �����
```
    docker build -t your_image_name:your_tag .
```
6. ��������� ��������� � ����� ����-������
```
    docker run -t your_image_name:your_tag ./app ./cases/new_test_case.txt
```
7. � ������ ���������� ����-����� � Makefile - ��������� ��� ����-�����
```
    docker run -t your_image_name:your_tag
```