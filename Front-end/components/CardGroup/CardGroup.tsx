import './style.css'
import React from 'react';
import axiosInstance from '../../axios';
import { useGroups } from '../../context/GroupsContext';
import { useNavigate } from 'react-router-dom';



interface CardGroupProps  {
  Id: number;
  Name: string;
} 

const CardGroup: React.FC<CardGroupProps> = ({Id, Name }) => {
    const { refreshGroups } = useGroups();
    let navigate = useNavigate()

    const DeleteGroup = async (Id: number) =>{
      try{
        const response = await axiosInstance.delete(
          `/api/group/deleteGroup/${Id}`
        );
        refreshGroups(); 
        console.log(response)
      }catch(error){
        console.log(error)
      }
    }
    const groupNavigate = (group_id: number) => {
      navigate(`/${group_id}/GroupData`)
    };
    return (
    <>
      <div id={Id.toString()}  className='ContainerGroupCard'>
        <div className='container-data'>
          <h2>{Name}</h2>
          <div className='Buttons-CardGroup'>
            <button className='open-button' onClick={() => groupNavigate(Id)}>Open</button>
            <button className='delete-button' onClick={() => DeleteGroup(Id)}>Delete</button>
          </div>
        </div>
      </div>
    </>
    );
};

export default CardGroup;
