import {useState, useEffect, useMemo, useReducer} from 'react'
import { BehaviorSubject } from 'rxjs';
const useReducerStream = (reducer, initState) => {
    const [state, dispatch] = useReducer(reducer, initState)
    const stream$ = useMemo(()=>new BehaviorSubject(state),[])
    useEffect(()=>{
        stream$.next(state)
    },[state])
    return [state, dispatch, stream$]
}
export default useReducerStream