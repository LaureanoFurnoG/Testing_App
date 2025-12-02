import './style.css'
import React, { useState } from 'react';
import axiosInstance from '../../axios';
import { Divider, Radio, Table } from 'antd';
import type { TableColumnsType, TableProps } from 'antd';

interface DataType {
  key: React.Key;
  name: string;
  delete: any;
}

interface TableUsers_Group  {
  URL: string;
} 

const TableUsers_Group: React.FC<TableUsers_Group> = ({ URL }) => {
    const [selectionType, setSelectionType] = useState<'checkbox' | 'radio'>('checkbox');
    const columns: TableColumnsType<DataType> = [
    {
        title: 'Name',
        dataIndex: 'name',
        render: (text: string) => <a>{text}</a>,
    },
    {
        title: 'Delete',
        dataIndex: 'delete',
    },
    ];

    const data: DataType[] = [
    {
        key: '1',
        name: 'John Brown',
        delete: <button>Delete</button>
    }
    ];

    return (
    <>
    <div>
        <Divider />
        <Table<DataType> 
            columns={columns}
            dataSource={data}
        />
    </div>
    </>
    );
};

export default TableUsers_Group;
