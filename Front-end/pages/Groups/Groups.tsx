import React, { useEffect, useState } from 'react';
import './style.css'
import CardGroup from '../../components/CardGroup/CardGroup'
import axiosInstance from '../../axios';
import { useAuth } from '../../auth/AuthProvider';
interface Group{
  id: number;
  name: string;
  path: string;
  subGroups: []
}
const GroupsManagement: React.FC = () => {
  const [Groups, setGroups] = useState<Group[]>([])
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      fetchAllGroups()
    }
  }, [isAuthenticated])

  const fetchAllGroups = async () => {
    try {
      const res = await axiosInstance.get("/api/group/showAllGroups")
      setGroups(res.data.Groups)
    } catch (e) {
      console.log(e)
    }
  }
  return (
    <div className='Container-groups-all'>
      <div className='createG-Header'>
        <h2>Groups</h2>
        <button>Create Group</button>
      </div>
      <div className='cards-groups'>
        {Groups.map(group =>(
          <CardGroup Id={group.id} Name={group.name}/>
        ))}
      
      </div>
    </div>
  );
};

export default GroupsManagement;
