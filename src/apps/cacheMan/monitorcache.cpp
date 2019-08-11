/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2018 Great Hill Corporation (http://greathill.com)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * This file was generated with makeClass. Edit only those parts of the code inside
 * of 'EXISTING_CODE' tags.
 */
#include <algorithm>
#include "monitorcache.h"

namespace qblocks {

//---------------------------------------------------------------------------
IMPLEMENT_NODE(CMonitorCache, CCache);

//---------------------------------------------------------------------------
static string_q nextMonitorcacheChunk(const string_q& fieldIn, const void *dataPtr);
static string_q nextMonitorcacheChunk_custom(const string_q& fieldIn, const void *dataPtr);

//---------------------------------------------------------------------------
void CMonitorCache::Format(ostream& ctx, const string_q& fmtIn, void *dataPtr) const {
    if (!m_showing)
        return;

    string_q fmt = (fmtIn.empty() ? expContext().fmtMap["monitorcache_fmt"] : fmtIn);
    if (fmt.empty()) {
        ctx << toJson();
        return;
    }

    // EXISTING_CODE
    // EXISTING_CODE

    while (!fmt.empty())
        ctx << getNextChunk(fmt, nextMonitorcacheChunk, this);
}

//---------------------------------------------------------------------------
string_q nextMonitorcacheChunk(const string_q& fieldIn, const void *dataPtr) {
    if (dataPtr)
        return reinterpret_cast<const CMonitorCache *>(dataPtr)->getValueByName(fieldIn);

    // EXISTING_CODE
    // EXISTING_CODE

    return fldNotFound(fieldIn);
}

//---------------------------------------------------------------------------------------------------
bool CMonitorCache::setValueByName(const string_q& fieldNameIn, const string_q& fieldValueIn) {
    string_q fieldName = fieldNameIn;
    string_q fieldValue = fieldValueIn;

    // EXISTING_CODE
    // EXISTING_CODE

    if (CCache::setValueByName(fieldName, fieldValue))
        return true;

    switch (tolower(fieldName[0])) {
        case 'm':
            if ( fieldName % "monitors" ) {
                CMonitorCacheItem item;
                string_q str = fieldValue;
                while (item.parseJson3(str)) {
                    monitors.push_back(item);
                    item = CMonitorCacheItem();  // reset
                }
                return true;
            }
            break;
        default:
            break;
    }
    return false;
}

//---------------------------------------------------------------------------------------------------
void CMonitorCache::finishParse() {
    // EXISTING_CODE
    // EXISTING_CODE
}

//---------------------------------------------------------------------------------------------------
bool CMonitorCache::Serialize(CArchive& archive) {

    if (archive.isWriting())
        return SerializeC(archive);

    // Always read the base class (it will handle its own backLevels if any, then
    // read this object's back level (if any) or the current version.
    CCache::Serialize(archive);
    if (readBackLevel(archive))
        return true;

    // EXISTING_CODE
    // EXISTING_CODE
    archive >> monitors;
    finishParse();
    return true;
}

//---------------------------------------------------------------------------------------------------
bool CMonitorCache::SerializeC(CArchive& archive) const {

    // Writing always write the latest version of the data
    CCache::SerializeC(archive);

    // EXISTING_CODE
    // EXISTING_CODE
    archive << monitors;

    return true;
}

//---------------------------------------------------------------------------
CArchive& operator>>(CArchive& archive, CMonitorCacheArray& array) {
    uint64_t count;
    archive >> count;
    array.resize(count);
    for (size_t i = 0 ; i < count ; i++) {
        ASSERT(i < array.capacity());
        array.at(i).Serialize(archive);
    }
    return archive;
}

//---------------------------------------------------------------------------
CArchive& operator<<(CArchive& archive, const CMonitorCacheArray& array) {
    uint64_t count = array.size();
    archive << count;
    for (size_t i = 0 ; i < array.size() ; i++)
        array[i].SerializeC(archive);
    return archive;
}

//---------------------------------------------------------------------------
void CMonitorCache::registerClass(void) {
    // only do this once
    if (HAS_FIELD(CMonitorCache, "schema")) return;

    CCache::registerClass();

    size_t fieldNum = 1000;
    ADD_FIELD(CMonitorCache, "schema",  T_NUMBER, ++fieldNum);
    ADD_FIELD(CMonitorCache, "deleted", T_BOOL,  ++fieldNum);
    ADD_FIELD(CMonitorCache, "showing", T_BOOL,  ++fieldNum);
    ADD_FIELD(CMonitorCache, "cname", T_TEXT,  ++fieldNum);
    ADD_FIELD(CMonitorCache, "monitors", T_OBJECT|TS_ARRAY, ++fieldNum);

    // Hide our internal fields, user can turn them on if they like
    HIDE_FIELD(CMonitorCache, "schema");
    HIDE_FIELD(CMonitorCache, "deleted");
    HIDE_FIELD(CMonitorCache, "showing");
    HIDE_FIELD(CMonitorCache, "cname");

    builtIns.push_back(_biCMonitorCache);

    // EXISTING_CODE
    // EXISTING_CODE
}

//---------------------------------------------------------------------------
string_q nextMonitorcacheChunk_custom(const string_q& fieldIn, const void *dataPtr) {
    const CMonitorCache *mon = reinterpret_cast<const CMonitorCache *>(dataPtr);
    if (mon) {
        switch (tolower(fieldIn[0])) {
            // EXISTING_CODE
            // EXISTING_CODE
            case 'p':
                // Display only the fields of this node, not it's parent type
                if ( fieldIn % "parsed" )
                    return nextBasenodeChunk(fieldIn, mon);
                // EXISTING_CODE
                // EXISTING_CODE
                break;

            default:
                break;
        }
    }

    return "";
}

//---------------------------------------------------------------------------
bool CMonitorCache::readBackLevel(CArchive& archive) {

    bool done = false;
    // EXISTING_CODE
    // EXISTING_CODE
    return done;
}

//---------------------------------------------------------------------------
CArchive& operator<<(CArchive& archive, const CMonitorCache& mon) {
    mon.SerializeC(archive);
    return archive;
}

//---------------------------------------------------------------------------
CArchive& operator>>(CArchive& archive, CMonitorCache& mon) {
    mon.Serialize(archive);
    return archive;
}

//---------------------------------------------------------------------------
string_q CMonitorCache::getValueByName(const string_q& fieldName) const {

    // Give customized code a chance to override first
    string_q ret = nextMonitorcacheChunk_custom(fieldName, this);
    if (!ret.empty())
        return ret;

    // Return field values
    switch (tolower(fieldName[0])) {
        case 'm':
            if ( fieldName % "monitors" || fieldName % "monitorsCnt" ) {
                size_t cnt = monitors.size();
                if (endsWith(toLower(fieldName), "cnt"))
                    return uint_2_Str(cnt);
                if (!cnt) return "";
                string_q retS;
                for (size_t i = 0 ; i < cnt ; i++) {
                    retS += monitors[i].Format();
                    retS += ((i < cnt - 1) ? ",\n" : "\n");
                }
                return retS;
            }
            break;
    }

    // EXISTING_CODE
    // EXISTING_CODE

    // Finally, give the parent class a chance
    return CCache::getValueByName(fieldName);
}

//-------------------------------------------------------------------------
ostream& operator<<(ostream& os, const CMonitorCache& item) {
    // EXISTING_CODE
    // EXISTING_CODE

    item.Format(os, "", nullptr);
    os << "\n";
    return os;
}

//---------------------------------------------------------------------------
const CBaseNode *CMonitorCache::getObjectAt(const string_q& fieldName, size_t index) const {
    if ( fieldName % "monitors" && index < monitors.size() )
        return &monitors[index];
    return NULL;
}

//---------------------------------------------------------------------------
// EXISTING_CODE
// EXISTING_CODE
}  // namespace qblocks
