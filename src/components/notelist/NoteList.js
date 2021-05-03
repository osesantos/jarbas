import React from "react";
import './NoteList.css';
import * as fs from 'fs';

const notesPath = '../../../notes/';
let notesList = [];

export default class NoteList extends React.Component {
    constructor(props) {
        super(props);
        notesList = this.getNotes();
    }

    getNotes() {
        return fs.readdirSync(notesPath);
    }

    render(){
        return (
            <div className="NoteList">
                <div className="SearchBar"></div>
                <div className="List">
                    {
                        notesList.map(note => <h3>note.name</h3>)
                    }
                </div>
            </div>
        );
    }
}