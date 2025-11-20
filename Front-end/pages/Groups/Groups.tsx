import React from 'react';
import './style.css'
import CardGroup from '../../components/CardGroup/CardGroup'
const GroupsManagement: React.FC = () => {
  return (
    <div className='Container-groups-all'>
      <div className='createG-Header'>
        <h2>Groups</h2>
        <button>Create Group</button>
      </div>
      <div className='cards-groups'>
        <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
        <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
        <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
        <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
        
      </div>
    </div>
  );
};

export default GroupsManagement;
