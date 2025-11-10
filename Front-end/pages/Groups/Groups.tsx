import React from 'react';
import './style.css'
import CardGroup from '../../components/CardGroup/CardGroup'
const GroupsManagement: React.FC = () => {
  return (
    <div className='Container-groups-all'>
      <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
      <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
      <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
      <CardGroup Id={2} Name='Hola' Endpoints='2' Members='2'/>
    </div>
  );
};

export default GroupsManagement;
