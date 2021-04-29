import api from './api'
import { Observable } from 'rxjs';

function rxapi(args){
    return Observable.create(observer => {
        api(args).then(data=>{
            observer.next(data)
            observer.complete()
        }).catch(err=>{
            observer.error(err)
        })
    })
}

export default rxapi