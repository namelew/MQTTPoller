import axios, { AxiosResponse } from 'axios';
import { IRequest, IResult } from 'interfaces/IExperiment';

export const api = axios.create({
    baseURL: 'http://localhost:8000/orquestrator'
});

export const startExperiment = async (expParam: IRequest) : Promise<IResult[]> => {
    const response = await axios.post<IResult[]>("/api/experiment/start", expParam);

    if (response.status !== 200) {
        console.log(response.statusText);
        return [];
    }

    return response.data;
};

export const deleteExperiment =async ( id: number ) : Promise<void> =>{
    const response = await axios.delete(`/api/experiment/${id}`);
    
    if (response.status !== 200) {
        console.log(response.statusText);
        return;
    }
}