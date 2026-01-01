import { Button, Divider, Table, TableColumnsType } from 'antd'
import './style.css'
import { CheckCircleOutlined, DeleteOutlined } from "@ant-design/icons";
import { useEffect, useState } from 'react';
import axiosInstance from '../../axios';
import { useGroups } from '../../context/GroupsContext';

interface DataType {
  Id: React.Key;
  Name: any;
}

const InvitationPanel: React.FC = () => {
    const [Invitations, setInv] = useState<DataType[]>([])
    const [Action, setAction] = useState<boolean>(false)
    const { refreshGroups } = useGroups();

    const getInvitations = async () =>{
        try{
            const response = await axiosInstance.get("/api/group/showInvitationGroups")
             const data = response.data.groupsInvitation.map((g: any) => ({
                Id: g.Id,
                Name: (
                    <div className="GroupData">
                        <p>{g.Name}</p>
                    </div>
                ),
            }));
            console.log(response)
            setInv(data)
        }catch(error){
            console.log(error)
        }
    }

    const decline = async (key: React.Key) => {
        try{
            
        }catch(error){
            console.log(error)
        }
    };

    const accept = async (key: React.Key) => {
        try{
           await axiosInstance.patch("/api/group/acceptInvitation", {GroupID: key})
           setAction(!Action)
           refreshGroups()
        }catch(error){
            console.log(error)
        }
    };

    const columns: TableColumnsType<DataType> = [
      {
        title: "Name",
        dataIndex: "Name",
      },
      {
        title: "Delete",
        dataIndex: "delete",
        render: (_, record) => (
          <div className="deleteCell">
            <Button type="primary" className="decline-btn" icon={<CheckCircleOutlined />}  onClick={() => decline(record.Id)}>
                Decline
            </Button>
          </div>
        ),
      },
      {
        title: "Accept",
        dataIndex: "accept",
        render: (_, record) => (
          <div className="acceptCell">
            <Button type="primary" className="accept-btn" icon={<CheckCircleOutlined />}  onClick={() => accept(record.Id)}>
                Accept
            </Button>
          </div>
        ),
      },
    ];

    useEffect(() =>{
        getInvitations()
    },[Action])
    return(
        <>
        <div className='Invitations-table-container'>
            <h1 className='text-title-inv'>Invitations</h1>
            <Divider className="divider" />

            <Table<DataType>
                columns={columns}
                dataSource={Invitations}
                showHeader={false}
            />
        </div>
        </>
    )
}

export default InvitationPanel